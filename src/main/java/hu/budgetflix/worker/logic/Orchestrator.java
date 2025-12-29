package hu.budgetflix.worker.logic;


import hu.budgetflix.worker.model.Status;
import hu.budgetflix.worker.model.Video;
import hu.budgetflix.worker.model.database.dao.MediaDao;
import hu.budgetflix.worker.view.Out;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class Orchestrator {

    private final FileMover mover;
    private final FfmpegRunner runner;
    private MediaDao dao;

    private final ExecutorService encoderExecutor = Executors.newSingleThreadExecutor();
    private final Object lock = new Object();
    private CompletableFuture<Void> tail = CompletableFuture.completedFuture(null);

    public Orchestrator(FileMover mover, FfmpegRunner runner, MediaDao dao) {
        this.mover = mover;
        this.runner = runner;
        this.dao = dao;
    }

    public void submit(Video video) {

        synchronized (lock) {
            tail = tail.thenRunAsync(() -> {
                        try {
                            enCode(video);
                        } catch (IOException e) {
                            throw new RuntimeException(e);
                        }
                    }, encoderExecutor)
                    .handle((v, ex) -> {
                        try {
                            if (ex != null) moveToError(video, ex);
                            else moveToDone(video);
                        } catch (IOException ioe) {
                            Out.log("post-processing failed: " + ioe);
                        }
                        return null;
                    });
        }
    }

    public void shutdownGracefully() {
        CompletableFuture<Void> drain;
        synchronized (lock) {
            drain = tail;
        }
        drain.join();
        encoderExecutor.shutdown();
        runner.shutdown();

    }

    private void moveToDone(Video video) throws IOException {
        video.setStatus(Status.DONE);
        dao.updateStatus(video);
        mover.moveProcessingToDone(video.getPath());

    }

    private void moveToError(Video video, Throwable exception) throws IOException {
        video.setStatus(Status.ERROR);
        dao.updateStatus(video);
        //Out.writeErrorLog(video,exception.toString());
        mover.moveProcessingToError(video.getPath());

    }

    private void enCode(Video video) throws IOException {
        video.setPath(mover.moveNewToProcessing(video.getPath()));


        Path outDir = Path.of("/srv/media/library/", video.getId().toString(), "hls");
        Files.createDirectories(outDir);


        List<String> cmd = buildFfmpegCmd(video, outDir);

        Out.log("start ffmpeg | " + video.toString());
        dao.addNewMedia(video, outDir);
        runner.start(cmd);
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
}
