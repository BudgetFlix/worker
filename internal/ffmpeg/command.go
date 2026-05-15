package ffmpeg

import (
	"path/filepath"
)

func BuildHLSCommand(
	input string,
	outputDir string,
) []string {

	segmentPattern := filepath.Join(
		outputDir,
		"seg_%03d.ts",
	)

	playlist := filepath.Join(
		outputDir,
		"index.m3u8",
	)

	return []string{
		"ffmpeg",

		"-threads", "1",

		"-y",

		"-i", input,

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
		segmentPattern,

		playlist,
	}
}