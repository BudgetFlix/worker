package hu.budgetflix.worker.logic;

import hu.budgetflix.worker.model.Status;


import java.io.IOException;
import java.nio.file.DirectoryStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.StandardCopyOption;


public class StatusExtension {

    public static Status getStatusExtension(Path path) {
        if (path == null) return null;

        String fileName = path.getFileName().toString();
        int lastDot = fileName.lastIndexOf(".");

        if (lastDot == -1) return null;

        String ext = fileName.substring(lastDot + 1).toUpperCase();

        try {
            return Status.valueOf(ext);
        } catch (IllegalArgumentException e) {
            return null;
        }
    }

    public static Path findMatching(Path videoPath) throws IOException {
        Path dir = videoPath.getParent();
        String baseName = videoPath.getFileName().toString();

        try (DirectoryStream<Path> stream =
                     Files.newDirectoryStream(dir, baseName + "*")) {

            for (Path entry : stream) {
                return entry; // első találat
            }
        }

        return null; // nincs találat
    }


    public static Path renameWithStatus(Path original, Status status) throws IOException, IOException {

        if (original == null || status == null) {
            throw new IllegalArgumentException("Path and Status cannot be null");
        }

        String fileName = original.getFileName().toString();

        Status existingStatus = getStatusExtension(original);
        if (existingStatus != null) {
            int lastDot = fileName.lastIndexOf(".");
            fileName = fileName.substring(0, lastDot);
        }

        String newFileName = fileName + "." + status.name().toLowerCase();
        Path target = original.resolveSibling(newFileName);

        return Files.move(
                original,
                target,
                StandardCopyOption.REPLACE_EXISTING,
                StandardCopyOption.ATOMIC_MOVE
        );

    }
}
