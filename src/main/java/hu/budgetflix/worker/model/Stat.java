package hu.budgetflix.worker.model;

import java.util.HashMap;

public class Stat {
    private String name;
    private MediaType type;
    private HashMap<Integer,String> videos;

    public Stat () {}

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public MediaType getType() {
        return type;
    }

    public void setType(MediaType type) {
        this.type = type;
    }

    public HashMap<Integer, String> getVideos() {
        return videos;
    }

    public void setVideos(HashMap<Integer, String> videos) {
        this.videos = videos;
    }
}
