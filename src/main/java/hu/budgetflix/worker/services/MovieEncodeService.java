package hu.budgetflix.worker.services;

import hu.budgetflix.worker.config.FfmpegConfig;
import hu.budgetflix.worker.config.WorkerConfig;
import hu.budgetflix.worker.logic.FfmpegRunner;
import hu.budgetflix.worker.logic.FileMover;
import hu.budgetflix.worker.logic.StatusExtension;
import hu.budgetflix.worker.model.JobResult;
import hu.budgetflix.worker.model.Stat;
import hu.budgetflix.worker.model.Status;
import hu.budgetflix.worker.model.database.JsonReader;
import hu.budgetflix.worker.model.database.dao.MediaDao;
import hu.budgetflix.worker.model.media.Movie;
import hu.budgetflix.worker.model.media.Video;
import hu.budgetflix.worker.view.Out;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;

public class MovieEncodeService implements EncodeService {

    private final MediaDao dao;
    private Movie movie;
    private boolean isShutDown = false;

    private final FfmpegRunner runner;

    public MovieEncodeService(MediaDao dao, FfmpegRunner runner) {
        this.dao = dao;
        this.runner = runner;
    }

    @Override
    public void buildUpStructure(Path directory) {
        try {

            Stat stat = JsonReader.jsonToObject(directory);

            String filename = stat.getVideos().get(1);

            Video video = new Video(StatusExtension.renameWithStatus(directory.resolve(filename), Status.READY));

            movie = new Movie(directory, stat, video);

            movie.setCurrentPath(FileMover.moveToProcessing(movie.getCurrentPath()));
            long id = dao.addNewMedia(movie);
            movie.setId(id);

            Path outPath = WorkerConfig.MOVIE_SOURCE.resolve(Long.toString(id));
            movie.setOutPath(outPath);

            dao.updateOutPatch(movie);
            Files.createDirectories(movie.getVideo().getOutPath());
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public void startEncode() {
        Video currentVideo = movie.getVideo();
        Status status = StatusExtension.getStatusExtension(currentVideo.getCurrentPath());

        if (status != Status.PROCESS && status != Status.READY) {
            return;
        }

        if(status == Status.PROCESS){
            //GC-torol minedt db-ben es localisan is
        }

        List<String> cmd = FfmpegConfig.buildFfmpegCmd(currentVideo);

        JobResult result;
        try {
            currentVideo.setCurrentPath(StatusExtension.renameWithStatus(currentVideo.getCurrentPath(),Status.PROCESS));
            result = runner.start(cmd, currentVideo.getFileName());
        } catch (IOException e) {
            throw new RuntimeException(e);
        }

        if (result != null) {

            if (result.success()) {
                try {
                    success(currentVideo);
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }
            } else {
                try {
                    failed(currentVideo,result);
                } catch (IOException e) {
                    throw new RuntimeException(e);
                }
            }
        }
        if(isShutDown){
            runner.shutdown();
        }
    }

    @Override
    public void shutDown () {
        isShutDown = true;
    }

    private void success (Video currentVideo) throws IOException {
        currentVideo.setCurrentPath(StatusExtension.renameWithStatus(currentVideo.getCurrentPath(),Status.DONE));
        dao.updateStatus(movie.getId(),Status.DONE);
        movie.setCurrentPath(FileMover.moveToDone(movie.getCurrentPath()));
    }

    private void failed (Video currentVideo,JobResult result) throws IOException {
        currentVideo.setCurrentPath(StatusExtension.renameWithStatus(currentVideo.getCurrentPath(),Status.ERROR));
        dao.updateStatus(movie.getId(),Status.ERROR);
        movie.setCurrentPath(FileMover.moveToError(movie.getCurrentPath()));
        Out.writeErrorLog(movie, result.errorTail());
    }


}
