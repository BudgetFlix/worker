package hu.budgetflix.worker.logic;

import hu.budgetflix.worker.model.JobResult;
import hu.budgetflix.worker.view.Out;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.nio.charset.StandardCharsets;
import java.util.ArrayDeque;
import java.util.Deque;
import java.util.List;
import java.util.concurrent.*;

public class FfmpegRunner {

    private Process currentFfmpeg;
    private final ExecutorService ioPool = Executors.newCachedThreadPool();

    public JobResult start(List<String> cmd,String filename) throws IOException {

        Out.log("ffmpeg started " + filename);

        ProcessBuilder pb = new ProcessBuilder(cmd);
        currentFfmpeg = pb.start();

        Deque<String> tail = new ArrayDeque<>(60);

        Future<?> stderrReader = ioPool.submit(() -> {
            try (BufferedReader br = new BufferedReader(
                    new InputStreamReader(currentFfmpeg.getErrorStream(), StandardCharsets.UTF_8))) {
                String line;
                while ((line = br.readLine()) != null) {
                    if (tail.size() == 60) tail.removeFirst();
                    tail.addLast(line);
                }
            } catch (IOException ignored) {}
        });

        Future<?> stdoutReader = ioPool.submit(() -> {
            try (InputStream in = currentFfmpeg.getInputStream()) {
                byte[] buf = new byte[8192];
                while (in.read(buf) != -1) {}
            } catch (IOException ignored) {}
        });

        int exit;
        try {
            exit = currentFfmpeg.waitFor();
            stderrReader.get(2, TimeUnit.SECONDS);
            stdoutReader.get(2, TimeUnit.SECONDS);

        } catch (InterruptedException | ExecutionException | TimeoutException e) {
            throw new RuntimeException("ffmpeg I/O handling failed", e);
        }

        if (exit != 0) {
            String errTail = String.join("\n", tail);
            return new JobResult(exit,errTail);
        }
        Out.log("ffmpeg finished OK");
        return new JobResult(exit,"Is OK");
    }

    public void shutdown() throws InterruptedException {
       ioPool.shutdown();
       ioPool.awaitTermination(Long.MAX_VALUE, TimeUnit.NANOSECONDS);
    }

}
