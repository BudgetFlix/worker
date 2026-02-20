package hu.budgetflix.worker;

import hu.budgetflix.worker.controller.EncodeController;
import hu.budgetflix.worker.factory.EncodeServiceFactory;
import hu.budgetflix.worker.logic.FfmpegRunner;
import hu.budgetflix.worker.logic.Observer;
import hu.budgetflix.worker.model.database.dao.MediaDao;
import hu.budgetflix.worker.model.database.dao.MediaDaoJdbc;
import hu.budgetflix.worker.model.database.manager.MediaDBManager;

public class Main {
    public static void main(String[] args) {

        FfmpegRunner runner = new FfmpegRunner();
        MediaDBManager manager = new MediaDBManager();
        MediaDao dao = new MediaDaoJdbc(manager.getDataSource());

        EncodeServiceFactory factory =
                new EncodeServiceFactory(dao, runner);

        EncodeController controller =
                new EncodeController(factory,runner);

        Observer observer = new Observer(controller);

        observer.finished()
                .thenRun(controller::shutdownGracefully)
                .join();
    }
}
