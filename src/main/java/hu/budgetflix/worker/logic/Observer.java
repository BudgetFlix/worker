package hu.budgetflix.worker.logic;

import hu.budgetflix.worker.config.WorkerConfig;

import java.io.File;
import java.io.IOException;
import java.nio.file.DirectoryStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.attribute.BasicFileAttributes;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

class FileState {
    long lastSize;
    long lastMtime;
    long stableSince;
}

public class Observer {
    Map<Path, FileState> states = new ConcurrentHashMap<>();
    Deque<Path> readyToEncode = new ArrayDeque<>();
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
        if(allFilesDownloaded()){
            watchingDownloaderFile.shutdown();
        }
        try (DirectoryStream<Path> stream = Files.newDirectoryStream(WorkerConfig.NEW_DIR)) {
            long now = System.currentTimeMillis();

            for (Path file : stream) {
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

                    if (now - state.stableSince >= 90_0000) {
                        readyToEncode.addLast(file);
                        states.remove(file);
                    }
                } else {
                    state.lastSize = size;
                    state.lastMtime = mtime;
                    state.stableSince = 0;
                }
            }

            states.keySet().removeIf(p -> !Files.exists(p));

        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    private boolean allFilesDownloaded() {
        return readyToEncode.size() == Objects.requireNonNull(new File(WorkerConfig.NEW_DIR.toUri()).listFiles()).length;
    }

    public Optional<Path> findNextInNew (){
        return Optional.ofNullable(readyToEncode.removeFirst());
    }

}
