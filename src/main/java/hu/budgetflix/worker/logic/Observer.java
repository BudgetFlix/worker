package hu.budgetflix.worker.logic;

import hu.budgetflix.worker.config.WorkerConfig;
import hu.budgetflix.worker.model.Status;
import hu.budgetflix.worker.model.Video;
import hu.budgetflix.worker.view.Out;


import java.io.IOException;
import java.nio.file.DirectoryStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.attribute.BasicFileAttributes;
import java.util.*;
import java.util.concurrent.*;

class FileState {
    long lastSize;
    long lastMtime;
    long stableSince;
}

public class Observer {
    Map<Path, FileState> states = new ConcurrentHashMap<>();
    Deque<Video> readyToEncode = new ConcurrentLinkedDeque<>();
    ScheduledExecutorService watchingDownloaderFile = Executors.newSingleThreadScheduledExecutor();


    public Observer () {
        setup();
    }

    private void setup () {
        watchingDownloaderFile.scheduleAtFixedRate(
                this::tick,0,5, TimeUnit.SECONDS);
    }

    public int readyCount() {
        return readyToEncode.size();
    }

    void tick() {

        Out.log("observer is running");

        try (DirectoryStream<Path> stream = Files.newDirectoryStream(WorkerConfig.NEW_DIR)) {
            long now = System.currentTimeMillis();

            for (Path file : stream) {
                Out.log(file.toString()); // debug
                if (!Files.isRegularFile(file)) continue;

                FileState state = states.computeIfAbsent(
                        file, p -> new FileState()
                );
                BasicFileAttributes attr = Files.readAttributes(file, BasicFileAttributes.class);
                long size = attr.size();
                long mtime = attr.lastModifiedTime().toMillis();

                if (size == state.lastSize && mtime == state.lastMtime) {
                    if (state.stableSince == 0) {
                        state.stableSince = now;
                    }

                    if (now - state.stableSince >= 10) {
                        readyToEncode.addLast(new Video(file,file.getFileName().toString(), Status.PROCESS));
                        states.remove(file);
                    }
                } else {
                    state.lastSize = size;
                    state.lastMtime = mtime;
                    state.stableSince = 0;
                }
            }
            states.keySet().removeIf(p -> !Files.exists(p));

            if(states.isEmpty()){
                watchingDownloaderFile.shutdown();
            }

        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }



    public Optional<Video> findNextInNew (){
        return Optional.ofNullable(readyToEncode.removeFirst());
    }

}
