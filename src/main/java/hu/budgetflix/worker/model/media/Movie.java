package hu.budgetflix.worker.model.media;

import hu.budgetflix.worker.model.Stat;

import java.nio.file.Path;

public  class Movie {
    private  long id;
    private Path currentPath;
    private  Path outPath;
    private final String name;
    private final Video video;
    private final Stat stat;
    private boolean canModifyTheOutPath = true;

    public Movie (Path currentPath, Stat stat, Video video) {
        this.currentPath = currentPath;
        this.name = stat.getName();
        this.video = video;
        this.stat = stat;
    }

    public long getId() {
        return id;
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

    public String getName() {
        return name;
    }

    public Video getVideo() {
        return video;
    }


    public void setOutPath(Path outPath) {
        if(canModifyTheOutPath){
            this.outPath = outPath;
            video.setOutPath(outPath.resolve("hls"));
            canModifyTheOutPath = false;
        }
    }

    public void setId(long id) {
        this.id = id;
    }
}
