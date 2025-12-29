package hu.budgetflix.worker.view;

import hu.budgetflix.worker.config.WorkerConfig;
import hu.budgetflix.worker.model.Video;

import java.io.BufferedWriter;
import java.io.FileWriter;
import java.io.IOException;
import java.util.Optional;

public class Out {

    public static void log (String msg) {
        System.out.printf("[%s] %s%n",
                java.time.LocalDateTime.now(), msg);
    }


    public static void writeErrorLog(Video currentProcessingVideo, String msg) {
        try (BufferedWriter writer = new BufferedWriter(new FileWriter(WorkerConfig.ERROR_LOG.toFile(),true))) {


            StringBuilder errorLog = new StringBuilder();

            errorLog.append("-".repeat(10));
            errorLog.append(System.lineSeparator());
            errorLog.append(java.time.LocalDateTime.now());
            errorLog.append(System.lineSeparator());
            errorLog.append("-".repeat(10));
            errorLog.append(System.lineSeparator());
            errorLog.append("Current video path: "); errorLog.append(currentProcessingVideo.getPath());
            errorLog.append(System.lineSeparator());
            errorLog.append("-".repeat(10));
            errorLog.append(System.lineSeparator());
            errorLog.append(System.lineSeparator());
            errorLog.append("-".repeat(10));
            errorLog.append(System.lineSeparator());
            errorLog.append(msg);
            errorLog.append(System.lineSeparator());
            errorLog.append("-".repeat(10));
            errorLog.append(System.lineSeparator().repeat(2));

            writer.write(errorLog.toString());
            log("text appended successfully to error log");
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }
}
