package hu.budgetflix.worker.model;

import java.util.HashMap;
import java.util.UUID;

public class Stat {
    private String name;
    private MediaType type;
    private UUID id;
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

    public UUID getId() {
        return id;
    }

    public void setId(UUID id) {
        this.id = id;
    }
}
