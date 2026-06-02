package ffmpeg

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func HLS(
	ctx context.Context,
	input string,
	outputDir string,
) error {
	if err := cleanHLSOutput(outputDir); err != nil {
		return err
	}

	args := BuildHLSCommand(
		input,
		outputDir,
	)

	cmd := exec.CommandContext(
		ctx,
		args[0],
		args[1:]...,
	)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	wg.Add(2)

	tail := make([]string, 0, 60)

	var tailMutex sync.Mutex

	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(stderr)

		for scanner.Scan() {

			line := scanner.Text()

			fmt.Println(line)

			tailMutex.Lock()

			if len(tail) >= 60 {
				tail = tail[1:]
			}

			tail = append(tail, line)

			tailMutex.Unlock()
		}
	}()

	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(stdout)

		for scanner.Scan() {
		}
	}()

	err = cmd.Wait()

	wg.Wait()

	if err != nil {

		tailMutex.Lock()
		output := strings.Join(tail, "\n")
		tailMutex.Unlock()

		return fmt.Errorf(
			"ffmpeg failed: %w\n%s",
			err,
			output,
		)
	}

	return nil
}

func cleanHLSOutput(
	outputDir string,
) error {
	// TODO: Read the HLS output path from job JSON and clean it immediately when crash recovery detects an interrupted job.
	playlist := filepath.Join(
		outputDir,
		"index.m3u8",
	)

	if err := removeIfExists(playlist); err != nil {
		return err
	}

	segments, err := filepath.Glob(
		filepath.Join(
			outputDir,
			"seg_*.ts",
		),
	)
	if err != nil {
		return err
	}

	for _, segment := range segments {
		if err := removeIfExists(segment); err != nil {
			return err
		}
	}

	return nil
}

func removeIfExists(
	path string,
) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
