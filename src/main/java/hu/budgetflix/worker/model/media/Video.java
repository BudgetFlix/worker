package hu.budgetflix.worker.model.media;

import hu.budgetflix.worker.model.Status;

import java.nio.file.Path;


public class Video {
    private Path currentPath;
    private  Path outPath;
    private Status status;

     public Video (Path path){
         this.currentPath = path;
         this.status = Status.PROCESS;
     }

    public Path getCurrentPath() {
        return currentPath;
    }

    public void setCurrentPath(Path currentPath) {
         this.currentPath = currentPath;
    }



    public Status getStatus() {
        return status;
    }

    public void setStatus(Status status) {
        this.status = status;
    }

    public Path getOutPath() {
        return outPath;
    }

    public void setOutPath(Path outPath) {
        this.outPath = outPath;
    }
}
