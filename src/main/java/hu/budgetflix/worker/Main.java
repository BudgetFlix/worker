package hu.budgetflix.worker;

import hu.budgetflix.worker.logic.FfmpegRunner;
import hu.budgetflix.worker.logic.FileMover;
import hu.budgetflix.worker.logic.Observer;
import hu.budgetflix.worker.logic.Orchestrator;

import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

public class Main {
    public static void main(String[] args) throws Exception {

        ScheduledExecutorService starter = Executors.newSingleThreadScheduledExecutor();
        System.out.println("hello worker");

        FileMover mover = new FileMover();
        FfmpegRunner runner = new FfmpegRunner();
        Observer observer = new Observer();

        Orchestrator orchestrator = new Orchestrator(mover,runner, observer);

        starter.scheduleAtFixedRate(() -> {
            if(observer.readyCount() > 0 ){
                try {
                    System.out.println("in progres");
                    orchestrator.runOnceUntilIdle();
                } catch (Exception e) {
                    throw new RuntimeException(e);
                }
                starter.shutdown();
            }
        },0,5, TimeUnit.SECONDS);



    }
}
