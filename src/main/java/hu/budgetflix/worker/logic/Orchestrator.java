package hu.budgetflix.worker.logic;

import hu.budgetflix.worker.model.JobResult;
import hu.budgetflix.worker.view.Out;

import java.io.IOException;
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

            Optional<Path> next = observer.findNextInNew();
            if (next.isPresent()) {
                Path moved;
                try {
                    moved = mover.moveNewToProcessing(next.get());
                    currentProcessingFile = moved;

                    List<String> cmd = buildFfmpegCmd(moved);

                    Out.log("start ffmpeg | " + currentProcessingFile);
                    currentJob = runner.start(cmd);

                } catch (Exception e) {
                    Out.log(e.getMessage());
                }

                continue;
            }

            break;
        }

    }

    private List<String> buildFfmpegCmd(Path in) {
        Path out = in.resolveSibling(in.getFileName().toString() + ".mp4");
        return List.of(
                "ffmpeg",
                "-y",
                "-i", "/srv/media/inbox/process/movie.mp4",

                "-map", "0:v:0",
                "-map", "0:a:0",

                "-c:v", "libx264",
                "-preset", "veryfast",
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
                "/srv/media/library/UUID/hls/seg_%03d.ts",

                "/srv/media/library/UUID/hls/index.m3u8"
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
