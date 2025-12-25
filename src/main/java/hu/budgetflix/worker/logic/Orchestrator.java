package hu.budgetflix.worker.logic;

import hu.budgetflix.worker.model.JobResult;
import hu.budgetflix.worker.model.Status;
import hu.budgetflix.worker.model.Video;
import hu.budgetflix.worker.model.database.dao.MediaDao;
import hu.budgetflix.worker.view.Out;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.CompletableFuture;

public class Orchestrator {

    private final FileMover mover;
    private final FfmpegRunner runner;
    private final Observer observer;
    private MediaDao dao;

    private CompletableFuture<JobResult> currentJob = null;
    private Optional<Video> currentProcessingVideo;

    public Orchestrator(FileMover mover, FfmpegRunner runner, Observer observer1, MediaDao dao) {
        this.mover = mover;
        this.runner = runner;
        this.observer = observer1;
        this.dao = dao;
    }

    public void runOnceUntilIdle() {

        while (true) {

            try {

                if (currentJob != null) {
                    Out.log("ffmpeg is in progress");
                    JobResult r = waitingForProcessingFinish();
                    onJobFinished(r);
                    continue;
                }

                currentProcessingVideo = observer.findNextInNew();
                if (currentProcessingVideo.isEmpty()) break;

                Video video = currentProcessingVideo.get();

                video.setPath(mover.moveNewToProcessing(video.getPath()));
                Out.log("success moving");

                Path outDir = Path.of("/srv/media/library/", video.getId().toString(), "hls");
                Files.createDirectories(outDir);


                List<String> cmd = buildFfmpegCmd(video, outDir);

                Out.log("start ffmpeg | " + currentProcessingVideo.get());
                dao.addNewMedia(currentProcessingVideo.get(), outDir);
                currentJob = runner.start(cmd);

            } catch (Exception e) {
                Out.log("in orchestrator: ");
                e.printStackTrace();
            }
        }


    }


    private List<String> buildFfmpegCmd(Video video, Path outDir) {
        return List.of(
                "ffmpeg",
                "-y",
                "-i", video.getPath().toString(),

                "-map", "0:v:0",
                "-map", "0:a:0",

                "-c:v", "libx264",
                "-preset", "fast",
                "-profile:v", "main",
                "-level", "4.0",
                "-pix_fmt", "yuv420p",
                "-crf", "20",

                "-c:a", "aac",
                "-b:a", "128k",
                "-ac", "2",

                "-f", "hls",
                "-hls_time", "6",
                "-hls_playlist_type", "vod",
                "-hls_segment_filename",
                outDir.resolve("seg_%03d.ts").toString(),
                outDir.resolve("index.m3u8").toString()
        );
    }

    private void onJobFinished(JobResult r) throws IOException {
        if (r.success()) {
            currentProcessingVideo.orElseThrow().setStatus(Status.DONE);
            dao.updateStatus(currentProcessingVideo.orElseThrow());

            mover.moveProcessingToDone(currentProcessingVideo.orElseThrow().getPath());
        } else {
            currentProcessingVideo.orElseThrow().setStatus(Status.ERROR);
            dao.updateStatus(currentProcessingVideo.orElseThrow());

            Out.writeErrorLog(currentProcessingVideo, r.exitCode(), r.errorTail());
            mover.moveProcessingToError(currentProcessingVideo.orElseThrow().getPath());
        }
        currentJob = null;
        currentProcessingVideo = Optional.empty();
    }

    private JobResult waitingForProcessingFinish() {
        return currentJob.join();
    }

}
