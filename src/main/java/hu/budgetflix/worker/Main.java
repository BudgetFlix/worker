package hu.budgetflix.worker;

import java.io.IOException;
import java.nio.file.*;
import java.time.LocalDateTime;
import java.util.stream.Stream;

public class Main {

        private static final Long SLEEP_MS = Long.parseLong(
                System.getenv().getOrDefault("WORKER_SLEEP_MS", "5000")
        );

        private static final Path INBOX_DIR = Paths.get(
                System.getenv().getOrDefault("INBOX_DIR", "/srv/media/inbox")
        );

    public static void main(String[] args) {

        while(true) {
            try {
                scanInbox();
                Thread.sleep(SLEEP_MS);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                log("Shot Down is requesting. By");
                return;
            } catch (Exception e) {
                log("Error in loop: " + e.getMessage());
                e.printStackTrace(System.out);
            }

        }



    }

    private static void scanInbox() throws IOException{
        if(!Files.exists(INBOX_DIR)){
            log("inbox is not exist" + INBOX_DIR);
            return;
        }
        if(!Files.isDirectory(INBOX_DIR)){
            log("inbox is not a directory" + INBOX_DIR);
            return;
        }
        try (Stream<Path> paths = Files.list(INBOX_DIR)){
            Long count = paths.filter(Files::isRegularFile).count();

            log("Inbox file sum: " + count);
        }
    }



    private static void log(String msg) {
        System.out.println(LocalDateTime.now() + " | " + msg);
    }
}
