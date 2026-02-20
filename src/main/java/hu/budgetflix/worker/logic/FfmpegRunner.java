package hu.budgetflix.worker.logic;

import hu.budgetflix.worker.model.JobResult;
import hu.budgetflix.worker.view.Out;
import hu.budgetflix.worker.view.StatusConsole;

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

    private final ExecutorService ioPool = Executors.newCachedThreadPool();

    public JobResult start(List<String> cmd,String filename) throws IOException {
        StatusConsole.setCurrentFile(filename);
        Out.log("ffmpeg started");

        ProcessBuilder pb = new ProcessBuilder(cmd);
        Process p = pb.start();

        Deque<String> tail = new ArrayDeque<>(60);

        Future<?> stderrReader = ioPool.submit(() -> {
            try (BufferedReader br = new BufferedReader(
                    new InputStreamReader(p.getErrorStream(), StandardCharsets.UTF_8))) {
                String line;
                while ((line = br.readLine()) != null) {
                    if (tail.size() == 60) tail.removeFirst();
                    tail.addLast(line);
                }
            } catch (IOException ignored) {}
        });

        Future<?> stdoutReader = ioPool.submit(() -> {
            try (InputStream in = p.getInputStream()) {
                byte[] buf = new byte[8192];
                while (in.read(buf) != -1) {}
            } catch (IOException ignored) {}
        });

        int exit;
        try {
            exit = p.waitFor();
            stderrReader.get(2, TimeUnit.SECONDS);
            stdoutReader.get(2, TimeUnit.SECONDS);

        } catch (InterruptedException | ExecutionException | TimeoutException e) {
            throw new RuntimeException("ffmpeg I/O handling failed", e);
        }

        if (exit != 0) {
            String errTail = String.join("\n", tail);
            StatusConsole.clearCurrentFile();
            return new JobResult(exit,errTail);
        }
        StatusConsole.clearCurrentFile();
        Out.log("ffmpeg finished OK");
        return new JobResult(exit,"Is OK");
    }

    public void shutdown() {
        ioPool.shutdown();
    }

}
