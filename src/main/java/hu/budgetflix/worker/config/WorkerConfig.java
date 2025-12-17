package hu.budgetflix.worker.config;

public class WorkerConfig {

    public static final String INBOX_DIR =
            System.getenv().getOrDefault("INBOX_DIR", "/srv/media/inbox");

    public static final String LIBRARY_DIR =
            System.getenv().getOrDefault("LIBRARY_DIR", "/srv/media/library");

    public static final String TEMP_DIR =
            System.getenv().getOrDefault("TEMP_DIR", "/srv/media/temp");
}
