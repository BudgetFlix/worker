package hu.budgetflix.worker.view;

public class Out {

    public static void log (String msg) {
        System.out.printf("[%s] %s%n",
                java.time.LocalDateTime.now(), msg);
    }

}
