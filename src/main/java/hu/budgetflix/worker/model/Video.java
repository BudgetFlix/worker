package hu.budgetflix.worker.model;

import java.nio.file.Path;
import java.util.UUID;

public class Video {
    private Path path;
    private final UUID id;
    private final String name;
    private Status status;

     public Video (Path path, String name, Status status){
         this.path = path;
         this.name = name;
         this.status = status;
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

    public String getName() {
        return name;
    }

    public Status getStatus() {
        return status;
    }

    public void setStatus(Status status) {
        this.status = status;
    }
}
