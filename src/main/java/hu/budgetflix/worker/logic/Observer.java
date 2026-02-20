package hu.budgetflix.worker.logic;

import hu.budgetflix.worker.config.WorkerConfig;
import hu.budgetflix.worker.controller.EncodeController;
import hu.budgetflix.worker.model.DirectoryState;
import hu.budgetflix.worker.view.Out;
import hu.budgetflix.worker.view.StatusConsole;

import java.io.IOException;
import java.nio.file.DirectoryStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.*;
import java.util.concurrent.*;

public class Observer {
    Map<Path, DirectoryState> states = new ConcurrentHashMap<>();
    ScheduledExecutorService watchingDownloaderFile = Executors.newSingleThreadScheduledExecutor();

    private final CompletableFuture<Void> finished =
            new CompletableFuture<>();
    private final EncodeController encodeController;


    public Observer (EncodeController encodeController) {
        this.encodeController = encodeController;
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


            for (Path file : stream) {
                file = file.toAbsolutePath().normalize();


                if (!Files.isDirectory(file)) continue;

                DirectoryState state = states.computeIfAbsent(
                        file, p -> new DirectoryState()
                );

                Out.log("stableSince: " + state.getStableSince()+ " " + file.getFileName());
                Out.log("submitted: " + state.isSubmitted()+ " " + file.getFileName());

                Path readyFile = isReadyToEncode(file);
                if(readyFile != null && !state.isSubmitted()) {
                    if (state.getStableSince() < 3) {
                        state.setStableSince(state.getStableSince() + 1);
                    }else{
                        state.setSubmitted(true);
                        Out.log(file + " is submitted");
                        encodeController.submit(file);
                    }
                }
            }
            states.keySet().removeIf(p -> !Files.exists(p));

            if (allstateIsSubmitted()
                    && encodeController.isIdle()
                    && !finished.isDone()) {

                Out.log("All encodes finished. Shutting down observer.");
                watchingDownloaderFile.shutdown();
                finished.complete(null);
            }

        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }

    private boolean allstateIsSubmitted() {
        return !states.isEmpty() &&
                states.values().stream()
                        .allMatch(DirectoryState::isSubmitted);
    }


    private Path isReadyToEncode(Path file) {
        if(!Files.exists(file)) {
            throw new RuntimeException("the given file is not exist");
        }

        String fileName = file.getFileName().toString();
        int lastDot = fileName.lastIndexOf('.');
        boolean hasExtension = lastDot > 0 && lastDot < fileName.length() - 1;

        if(hasExtension &&  fileName.endsWith(".part")){
            return null;
        }
        return file;
    }
}
