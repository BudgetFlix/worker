package hu.budgetflix.worker.controller;

import hu.budgetflix.worker.logic.FfmpegRunner;
import hu.budgetflix.worker.model.database.JsonReader;
import hu.budgetflix.worker.factory.EncodeServiceFactory;
import hu.budgetflix.worker.model.MediaType;
import hu.budgetflix.worker.services.EncodeService;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;

public class EncodeController {

    private final ExecutorService executor =
            Executors.newSingleThreadExecutor();

    private final EncodeServiceFactory factory;
    private FfmpegRunner runner;

    public EncodeController(EncodeServiceFactory factory, FfmpegRunner runner) {
        this.factory = factory;
        this.runner = runner;
    }

    public void submit(Path directory) {
        executor.submit(() -> process(directory));
    }

    private void process(Path directory) {
        try {
            validate(directory);

            MediaType type = JsonReader.extractTheType(directory);
            EncodeService service = factory.create(type);

            service.buildUpStructure(directory);
            service.startEncode();

        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    private void validate(Path directory) {
        if (!Files.exists(directory)) {
            throw new IllegalArgumentException(
                    "Directory does not exist: " + directory);
        }
    }

    public void shutdownGracefully() {
        executor.shutdown();
        try {
            executor.awaitTermination(1, TimeUnit.HOURS);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
        runner.shutdown();
    }
}

