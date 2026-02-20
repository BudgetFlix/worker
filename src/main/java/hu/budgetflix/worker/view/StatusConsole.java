package hu.budgetflix.worker.view;


public class StatusConsole {

    private static volatile boolean watching = false;
    private static volatile String currentFile = "-";

    public static void setWatching(boolean value) {
        watching = value;
        render();
    }

    public static void setCurrentFile(String file) {
        currentFile = file;
        render();
    }

    public static void clearCurrentFile() {
        currentFile = "-";
        render();
    }

    private static synchronized void render() {

        String statusLine = String.format(
                "Observer: %-5s | FFmpeg: %-30s",
                watching, currentFile
        );

        // kurzor mentése
        System.out.print("\033[s");

        // menj le a legalsó sorra (egyszerűbb: egy sort fel az aktuális végéről)
        System.out.print("\033[999B");   // nagyon le
        System.out.print("\033[2K");     // sor törlése
        System.out.print("\r" + statusLine);

        // kurzor vissza
        System.out.print("\033[u");

        System.out.flush();
    }
}
