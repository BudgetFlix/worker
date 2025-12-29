package hu.budgetflix.worker.logic;

import hu.budgetflix.worker.config.WorkerConfig;
import hu.budgetflix.worker.model.Status;
import hu.budgetflix.worker.model.Video;
import hu.budgetflix.worker.view.Out;


import javax.swing.plaf.nimbus.State;
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
    boolean submitted = false;
}

public class Observer {
    Map<Path, FileState> states = new ConcurrentHashMap<>();
    ScheduledExecutorService watchingDownloaderFile = Executors.newSingleThreadScheduledExecutor();
    private final Orchestrator orchestrator;
    private final CompletableFuture<Void> finished =
            new CompletableFuture<>();


    public Observer (Orchestrator orchestrator1) {
        this.orchestrator = orchestrator1;
        setup();
    }

    private void setup () {
        watchingDownloaderFile.scheduleAtFixedRate(
                this::tick,0,5, TimeUnit.SECONDS);
    }

    public CompletableFuture<Void> finished() {
        return finished;
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

                    if (now - state.stableSince >= 20_000 && !state.submitted) {
                        orchestrator.submit(new Video(file,file.getFileName().toString(), Status.PROCESS));
                        state.submitted = true;
                    }
                } else {
                    state.lastSize = size;
                    state.lastMtime = mtime;
                    state.stableSince = 0;
                }
            }
            states.keySet().removeIf(p -> !Files.exists(p));

            if(allstateIsSubmitted() && !finished.isDone()){
                watchingDownloaderFile.shutdown();
                finished.complete(null);
            }

        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    private boolean allstateIsSubmitted () {
        return states.values().stream().allMatch(s -> s.submitted);
    }
}
