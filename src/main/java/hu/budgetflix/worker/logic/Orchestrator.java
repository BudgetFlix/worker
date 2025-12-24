package hu.budgetflix.worker.logic;

import hu.budgetflix.worker.model.JobResult;
import hu.budgetflix.worker.model.Video;
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

    private CompletableFuture<JobResult> currentJob = null;
    private Optional<Video> currentProcessingVideo;

    public Orchestrator(FileMover mover, FfmpegRunner runner, Observer observer1) {
        this.mover = mover;
        this.runner = runner;
        this.observer = observer1;
    }

    public void runOnceUntilIdle() throws Exception {

        while (true) {


            if (currentJob != null) {
                Out.log("ffmpeg is in progress");
                JobResult r = waitingForProcessingFinish();
                onJobFinished(r);
                continue;
            }

            currentProcessingVideo = observer.findNextInNew();
            if (currentProcessingVideo.isEmpty()) break;

            try {
                Video video = currentProcessingVideo.get();

                video.setPath(mover.moveNewToProcessing(video.getPath()));
                Out.log("success moving");

                Path outDir = Path.of("/srv/media/library/", video.getId().toString(), "/hls");
                Files.createDirectories(outDir);

                List<String> cmd = buildFfmpegCmd(video,outDir);

                Out.log("start ffmpeg | " + currentProcessingVideo.get().toString());
                currentJob = runner.start(cmd);

            } catch (Exception e) {
                Out.log("in orchestrator" + e + e.getMessage());
            }
        }


    }


    private List<String> buildFfmpegCmd(Video video,Path outDir) {
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
            mover.moveProcessingToDone(currentProcessingVideo.orElseThrow().getPath());
        } else {
            //mover.writeErrorLog(currentProcessingFile,r.exitCode(),r.exitCode());
            mover.moveProcessingToError(currentProcessingVideo.orElseThrow().getPath());
        }
        currentJob = null;
        currentProcessingVideo = Optional.empty();
    }

    private JobResult waitingForProcessingFinish() {
        return currentJob.join();
    }

}
