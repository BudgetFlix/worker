package hu.budgetflix.worker.model.media;

import hu.budgetflix.worker.model.Stat;
import hu.budgetflix.worker.model.Status;

import java.nio.file.Path;

public  class Movie {
    private  long id;
    private Path currentPath;
    private  Path outPath;
    private final String name;
    private final Video video;
    private Status status;
    private final Stat stat;
    private boolean canModifyTheOutPath = true;

    public Movie (Path currentPath, Stat stat) {
        this.currentPath = currentPath;
        this.name = stat.getName();
        this.video = new Video(currentPath.resolve(stat.getVideos().get(1)));
        this.stat = stat;
        this.status = Status.PROCESS;
    }

    public long getId() {
        return id;
    }

    public Path getCurrentPath() {
        return currentPath;
    }

    public void setCurrentPath(Path currentPath) {
        this.currentPath = currentPath;
        video.setCurrentPath(currentPath.resolve(stat.getVideos().get(1)));
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

    public Status getStatus() {
        return status;
    }

    public void setStatus(Status status) {
        this.status = status;
    }

    public Stat getStat() {
        return stat;
    }

    public void setOutPath(Path outPath) {
        if(canModifyTheOutPath){
            this.outPath = outPath;
            video.setOutPath(outPath.resolve(stat.getVideos().get(1) + "/hls"));
        }
    }

    public void setId(long id) {
        this.id = id;
    }
}
