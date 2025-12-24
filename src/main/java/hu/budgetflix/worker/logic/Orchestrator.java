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
    private Path currentProcessingFile = null;

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

            Optional<Video> next = observer.findNextInNew();
            if (next.isPresent()) {
                Path moved;
                try {
                    moved = mover.moveNewToProcessing(next.get().getPath());
                    Out.log("success moving");
                    currentProcessingFile = moved;

                    List<String> cmd = buildFfmpegCmd(next);

                    Out.log("start ffmpeg | " + currentProcessingFile);

                    Path out = Path.of("/srv/media/library/" + next.get().getId() + "/hls/");
                    Files.createDirectories(out);
                    currentJob = runner.start(cmd);

                } catch (Exception e) {
                    Out.log(e.getMessage());
                }

                continue;
            }

            break;
        }

    }

    private List<String> buildFfmpegCmd(Optional<Video> video) {
        if(video.isEmpty()) return List.of();
        Path out = Path.of("/srv/media/library/" + video.get().getId() + "/hls/");
        return List.of(
                "ffmpeg",
                "-y",
                "-i", video.get().getPath().toString(),

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
                out + "seg_%03d.ts",

                out + "index.m3u8"
        );
    }

    private void onJobFinished(JobResult r) throws IOException {
        if (r.success()) {
            mover.moveProcessingToDone(currentProcessingFile);
        } else {
            //mover.writeErrorLog(currentProcessingFile,r.exitCode(),r.exitCode());
            mover.moveProcessingToError(currentProcessingFile);
        }
        currentJob = null;
        currentProcessingFile = null;
    }

    private JobResult waitingForProcessingFinish() {
        return currentJob.join();
    }

}
