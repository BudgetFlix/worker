package hu.budgetflix.worker.services;

import hu.budgetflix.worker.config.FfmpegConfig;
import hu.budgetflix.worker.config.WorkerConfig;
import hu.budgetflix.worker.logic.FfmpegRunner;
import hu.budgetflix.worker.logic.FileMover;
import hu.budgetflix.worker.model.JobResult;
import hu.budgetflix.worker.model.Stat;
import hu.budgetflix.worker.model.Status;
import hu.budgetflix.worker.model.database.JsonReader;
import hu.budgetflix.worker.model.database.dao.MediaDao;
import hu.budgetflix.worker.model.media.Movie;
import hu.budgetflix.worker.view.Out;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;

public class MovieEncodeService implements EncodeService {


    private final MediaDao dao;
    private Movie movie;

    private FfmpegRunner runner;

    public MovieEncodeService(MediaDao dao, FfmpegRunner runner) {
        this.dao = dao;
        this.runner = runner;
    }

    @Override
    public void buildUpStructure(Path directory) {
        try {

        Stat stat = JsonReader.jsonToObject(directory);
        movie = new Movie(directory, stat);
        long id = dao.addNewMedia(movie);
        Path outPath = WorkerConfig.MOVIE_SOURCE.resolve(Long.toString(id));
        movie.setId(id);
        movie.setOutPath(outPath);
        dao.updatePatch(movie);
        Files.createDirectories(movie.getVideo().getOutPath());

        movie.setCurrentPath(FileMover.moveNewToProcessing(movie.getCurrentPath()));

        } catch (Exception e) {
            throw new RuntimeException(e);
        }


    }

    @Override
    public void startEncode() {
        List<String> cmd = FfmpegConfig.buildFfmpegCmd(movie.getVideo());

        JobResult result;
        try {
            result = runner.start(cmd,movie.getVideo().getCurrentPath().getFileName().toString());
        } catch (IOException e) {
            throw new RuntimeException(e);
        }

        if(result != null) {

        if (result.success()){ // athidalni statuszt a movi-re is
                try {
                    movie.setStatus(Status.DONE);
                    dao.updateStatus(movie);
                    movie.setCurrentPath( FileMover.moveProcessingToDone(movie.getCurrentPath()));
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }
            }else {
                try {
                    movie.setStatus(Status.ERROR);
                    dao.updateStatus(movie);
                    movie.setCurrentPath( FileMover.moveProcessingToError(movie.getCurrentPath()));
                    Out.writeErrorLog(movie,result.errorTail());
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }
            }
        }
    }

}
