package hu.budgetflix.worker.logic;

import hu.budgetflix.worker.config.WorkerConfig;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;

public class FileMover {
    public Path moveNewToProcessing(Path path) throws IOException {
        return Files.move(path, WorkerConfig.PROCESS_DIR.resolve(path.getFileName()),
                StandardCopyOption.ATOMIC_MOVE,
                StandardCopyOption.REPLACE_EXISTING);
    }

    public void moveProcessingToDone(Path currentProcessingFile) throws IOException {
        Files.move(currentProcessingFile, WorkerConfig.DONE_DIR.resolve(currentProcessingFile.getFileName()),
                StandardCopyOption.ATOMIC_MOVE,
                StandardCopyOption.REPLACE_EXISTING);
    }

    public void writeErrorLog(Path currentProcessingFile, int i, int i1) {

    }

    public void moveProcessingToError(Path currentProcessingFile) throws IOException {
        Files.move(currentProcessingFile, WorkerConfig.ERROR_DIR.resolve(currentProcessingFile.getFileName()),
                StandardCopyOption.ATOMIC_MOVE,
                StandardCopyOption.REPLACE_EXISTING);
    }
}
