package hu.budgetflix.worker.model.database.dao;

import hu.budgetflix.worker.model.Video;

import java.nio.file.Path;

public interface MediaDao {
    void addNewMedia (Video vide, Path out);
    void updateStatus(Video vide);
}
