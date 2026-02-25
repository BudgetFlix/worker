package hu.budgetflix.worker.model.media;

import hu.budgetflix.worker.model.Status;

import java.nio.file.Path;


public class Video {
    private Path currentPath;
    private  Path outPath;
    private String fileName;

     public Video (Path path){
         this.currentPath = path;
         this.fileName = path.getFileName().toString();
     }

    public Path getCurrentPath() {
        return currentPath;
    }

    public void setCurrentPath(Path currentPath) {
         this.currentPath = currentPath;
    }


    public Path getOutPath() {
        return outPath;
    }

    public void setOutPath(Path outPath) {
        this.outPath = outPath;
    }

    public String getFileName() {
        return fileName;
    }
}
