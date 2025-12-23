package hu.budgetflix.worker.logic;

import hu.budgetflix.worker.model.JobResult;

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
                JobResult r = waitingForProcessingFinish();
                onJobFinished(r);
                continue;
            }

            Optional<Path> next = observer.findNextInNew();
            if (next.isPresent()) {
                Path moved = mover.moveNewToProcessing(next.get());
                currentProcessingFile = moved;

                List<String> cmd = buildFfmpegCmd(moved);
                currentJob = runner.start(cmd);
                continue;
            }

            break;
        }

    }

    private List<String> buildFfmpegCmd(Path in) {
        Path out = in.resolveSibling(in.getFileName().toString() + ".mp4");
        return List.of("ffmpeg", "-y", "-i", in.toString(), out.toString());
    }

    private void onJobFinished(JobResult r) throws IOException {
        if(r.success()){
            mover.moveProcessingToDone(currentProcessingFile);
        } else {
            mover.writeErrorLog(currentProcessingFile,r.exitCode(),r.exitCode());
            mover.moveProcessingToError(currentProcessingFile);
        }
        currentJob = null;
        currentProcessingFile = null;
    }

    private JobResult waitingForProcessingFinish() {
        return currentJob.join();
    }

}
