package hu.budgetflix.worker.model;

import java.nio.file.Path;
import java.util.UUID;

public class Video {
    private final Path path;
    private final UUID id;

     public Video (Path path){
         this.path = path;
         this.id = UUID.randomUUID();
     }

    public Path getPath() {
        return path;
    }

    public UUID getId() {
        return id;
    }
}
