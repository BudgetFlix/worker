package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"worker/internal/job"
)

func Job(job *job.MediaJob) {
	log.Println("========================================")
	log.Println("MEDIA JOB")
	log.Println("========================================")

	log.Printf("ID:        %s\n", job.ID)
	log.Printf("Type:      %s\n", job.Type)
	log.Printf("Path:      %s\n", job.Path)
	log.Printf("Items:     %d\n", len(job.Items))

	log.Println("----------------------------------------")

	logRealPaths(job)

	log.Println("========================================")
	fmt.Println()
}

func logRealPaths(job *job.MediaJob) {
	log.Println("========================================")
	log.Println("REAL FILE PATHS")
	log.Println("========================================")

	for i, item := range job.Items {

		pattern := filepath.Join(job.Path, item.FileName+"*")

		matches, err := filepath.Glob(pattern)
		if err != nil {
			log.Printf("ITEM #%d -> glob error: %v\n", i+1, err)
			continue
		}

		if len(matches) == 0 {
			log.Printf("ITEM #%d -> no file found for: %s\n", i+1, item.FileName)
			continue
		}

		log.Printf("ITEM #%d\n", i+1)
		log.Printf("  FileName: %s\n", item.FileName)

		for _, match := range matches {

			absPath, err := filepath.Abs(match)
			if err != nil {
				absPath = match
			}

			info, err := os.Stat(match)
			if err != nil {
				log.Printf("  Path: %s (stat error: %v)\n", absPath, err)
				continue
			}

			log.Printf("  Found: %s\n", absPath)

			if info.IsDir() {
				log.Printf("  Type: DIRECTORY\n")
			} else {
				log.Printf("  Type: FILE\n")
			}
		}

		log.Println("----------------------------------------")
	}

	log.Println("========================================")
}

func Loging(msg string) {

	symbol := "="
	for range len(msg) {
		fmt.Print(symbol)
	}

	fmt.Println()
	fmt.Println(msg)

	for range len(msg) {
		fmt.Print(symbol)
	}
}
