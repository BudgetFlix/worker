package hu.budgetflix.worker.model;

import java.nio.file.Path;
import java.util.UUID;

public class Video {
    private Path path;
    private final UUID id;

     public Video (Path path){
         this.path = path;
         this.id = UUID.randomUUID();
     }

    public Path getPath() {
        return path;
    }

    public void setPath (Path path) {
         this.path = path;
    }

    public UUID getId() {
        return id;
    }
}
