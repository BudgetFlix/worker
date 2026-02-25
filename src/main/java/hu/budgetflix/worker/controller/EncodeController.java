package hu.budgetflix.worker.controller;

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
    private volatile boolean shuttingDown = false;

    public EncodeController(EncodeServiceFactory factory) {
        this.factory = factory;
    }

    public void submit(Path directory) {
        executor.submit(() -> {
            if(shuttingDown) return;
            process(directory);
        });
    }

    private void process(Path directory) {
        try {
            validate(directory);

            MediaType type = JsonReader.extractTheType(directory);
            EncodeService service = factory.create(type);

            service.buildUpStructure(directory);
            service.startEncode();
            if(shuttingDown){
                service.shutDown();
            }

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


    public void shutdown() {
        shuttingDown = true;
        executor.shutdown();
        try {
            executor.awaitTermination(Long.MAX_VALUE, TimeUnit.NANOSECONDS);
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }
}

