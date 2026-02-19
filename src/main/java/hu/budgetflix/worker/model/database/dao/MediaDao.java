package hu.budgetflix.worker.model.database.dao;

import hu.budgetflix.worker.model.media.Movie;

public interface MediaDao {
    Long addNewMedia (Movie movie);
    void updatePatch (Movie movie);
    void updateStatus(Movie movie);
}
