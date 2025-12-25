package hu.budgetflix.worker;

import hu.budgetflix.worker.logic.FfmpegRunner;
import hu.budgetflix.worker.logic.FileMover;
import hu.budgetflix.worker.logic.Observer;
import hu.budgetflix.worker.logic.Orchestrator;
import hu.budgetflix.worker.model.database.dao.MediaDao;
import hu.budgetflix.worker.model.database.dao.MediaDaoJdbc;
import hu.budgetflix.worker.model.database.manager.MediaDBManager;
import hu.budgetflix.worker.view.Out;

import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

public class Main {
    public static void main(String[] args) {

        ScheduledExecutorService starter = Executors.newSingleThreadScheduledExecutor();
        System.out.println("hello worker");

        FileMover mover = new FileMover();
        FfmpegRunner runner = new FfmpegRunner();
        Observer observer = new Observer();

        MediaDBManager manager = new MediaDBManager();
        MediaDao dao = new MediaDaoJdbc(manager.getDataSource());

        Orchestrator orchestrator = new Orchestrator(mover, runner, observer,dao);

        starter.scheduleAtFixedRate(() -> {

            System.out.println(observer.readyCount());

            if (observer.readyCount() > 0) {
                try {
                    System.out.println("in progres");
                    orchestrator.runOnceUntilIdle();
                } catch (Exception e) {
                    Out.log(e.toString());
                }
                starter.shutdown();
            }
        }, 0, 5, TimeUnit.SECONDS);


    }
}
