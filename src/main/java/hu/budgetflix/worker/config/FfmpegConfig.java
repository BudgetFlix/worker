package hu.budgetflix.worker.config;

import hu.budgetflix.worker.model.Stat;
import hu.budgetflix.worker.model.media.Video;

import java.util.List;

public class FfmpegConfig {

    public static List<String> buildFfmpegCmd(Video video) {
        return List.of(
                "ffmpeg",
                "-threads", "1",
                "-y",
                "-i", video.getCurrentPath().toString(),

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
                video.getOutPath().resolve("seg_%03d.ts").toString(),
                video.getOutPath().resolve("index.m3u8").toString()
        );
    }
}
