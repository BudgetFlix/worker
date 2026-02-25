package hu.budgetflix.worker.model.database.dao;

import hu.budgetflix.worker.model.Status;
import hu.budgetflix.worker.model.media.Movie;

public interface MediaDao {
    Long addNewMedia (Movie movie);
    void updateOutPatch (Movie movie);
    void updateStatus(Long id, Status status);
}
