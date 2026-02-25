package hu.budgetflix.worker.services;

import hu.budgetflix.worker.logic.FfmpegRunner;
import hu.budgetflix.worker.model.database.dao.MediaDao;

import java.nio.file.Path;

public class SeriesEncodeService implements EncodeService {
    public SeriesEncodeService(MediaDao dao, FfmpegRunner runner) {

    }

    @Override
    public void buildUpStructure(Path directory) {

    }

    @Override
    public void startEncode() {

    }

    @Override
    public void shutDown() {

    }
}
