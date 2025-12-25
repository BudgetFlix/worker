package hu.budgetflix.worker.model.database.manager;

import hu.budgetflix.worker.config.WorkerConfig;
import org.sqlite.SQLiteConfig;
import org.sqlite.SQLiteDataSource;

import javax.sql.DataSource;

public class MediaDBManager {

    private final SQLiteDataSource dataSource;

    public MediaDBManager() {
        SQLiteConfig config = new SQLiteConfig();
        config.setJournalMode(SQLiteConfig.JournalMode.WAL);
        config.setSynchronous(SQLiteConfig.SynchronousMode.NORMAL);
        config.deferForeignKeys(true);
        config.setBusyTimeout(10_000);

        dataSource = new SQLiteDataSource();
        dataSource.setConfig(config);
        dataSource.setUrl("jdbc:sqlite:" + WorkerConfig.DATA_BASE);
    }

    public DataSource getDataSource() {
        return dataSource;
    }
}

