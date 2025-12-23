package hu.budgetflix.worker.config;

import java.nio.file.Path;

public class WorkerConfig {

    public final static Path NEW_DIR =
            Path.of(System.getenv().getOrDefault("NEW_DIR", "/srv/media/inbox/new"));

    public final static Path PROCESS_DIR =
            Path.of(System.getenv().getOrDefault("PROCESS_DIR", "/srv/media/inbox/process"));

    public final static Path DONE_DIR =
            Path.of(System.getenv().getOrDefault("DONE_DIR", "/srv/media/inbox/done"));

    public final static Path ERROR_DIR =
            Path.of(System.getenv().getOrDefault("ERROR_DIR", "/srv/media/inbox/error"));
}
