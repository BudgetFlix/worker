package hu.budgetflix.worker.logic;

import hu.budgetflix.worker.model.JobResult;

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

    public CompletableFuture<JobResult> start(List<String> cmd) throws IOException {
        ProcessBuilder pb = new ProcessBuilder(cmd);
        pb.redirectErrorStream(false);
        Process p = pb.start();

        int tailLines = 80;
        Deque<String> tail = new ArrayDeque<>(tailLines);

        Future<?> stderrReader = ioPool.submit(() -> {
            try (BufferedReader br = new BufferedReader(
                    new InputStreamReader(p.getErrorStream(), StandardCharsets.UTF_8))) {
                String line;
                while ((line = br.readLine()) != null){
                    if(tail.size() == tailLines) tail.removeFirst();
                    tail.addLast(line);
                }
            } catch (IOException ignored) {
            }
        });

        Future<?> stdoutReader = ioPool.submit(() -> {
            try(InputStream in = p.getInputStream()){
                byte[] buf = new byte[8192];
                while (in.read(buf) != -1){/*none*/}
            }catch (IOException ignored){}
        });

        CompletableFuture<JobResult> finished = p.onExit().thenApply(process -> {
            int exit = process.exitValue();

            try {stdoutReader.get(2, TimeUnit.SECONDS);}catch (Exception ignored){}
            try {stderrReader.get(2, TimeUnit.SECONDS);}catch (Exception ignored){}

            String errTail = String.join("\n", tail);
            return new JobResult(exit,errTail);
        });
        return finished;
    }
}
