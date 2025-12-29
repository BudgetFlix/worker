package hu.budgetflix.worker;

import hu.budgetflix.worker.logic.FfmpegRunner;
import hu.budgetflix.worker.logic.FileMover;
import hu.budgetflix.worker.logic.Observer;
import hu.budgetflix.worker.logic.Orchestrator;
import hu.budgetflix.worker.model.database.dao.MediaDao;
import hu.budgetflix.worker.model.database.dao.MediaDaoJdbc;
import hu.budgetflix.worker.model.database.manager.MediaDBManager;


public class Main {
    public static void main(String[] args) {

        System.out.println("hello worker");

        FileMover mover = new FileMover();
        FfmpegRunner runner = new FfmpegRunner();

        MediaDBManager manager = new MediaDBManager();
        MediaDao dao = new MediaDaoJdbc(manager.getDataSource());

        Orchestrator orchestrator = new Orchestrator(mover, runner,dao);

        Observer observer = new Observer(orchestrator);

        observer.finished().join();
        orchestrator.shutdownGracefully();
    }
}
