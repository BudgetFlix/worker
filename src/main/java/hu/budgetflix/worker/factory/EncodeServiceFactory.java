package hu.budgetflix.worker.factory;

import hu.budgetflix.worker.logic.FfmpegRunner;
import hu.budgetflix.worker.model.MediaType;
import hu.budgetflix.worker.model.database.dao.MediaDao;
import hu.budgetflix.worker.services.EncodeService;
import hu.budgetflix.worker.services.MovieEncodeService;
import hu.budgetflix.worker.services.SeriesEncodeService;



public class EncodeServiceFactory {

    private final MediaDao dao;
    private final FfmpegRunner runner;

    public EncodeServiceFactory(MediaDao dao, FfmpegRunner runner) {
        this.dao = dao;
        this.runner = runner;
    }

    public EncodeService create(MediaType type) {
        return switch (type) {
            case MOVIE -> new MovieEncodeService(dao, runner);
            case SERIES -> new SeriesEncodeService(dao, runner);
        };
    }
}
