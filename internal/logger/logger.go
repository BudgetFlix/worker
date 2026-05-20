package logger

import (
	"fmt"
	"log"

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

	for i, item := range job.Items {
		log.Printf("ITEM #%d\n", i+1)
		log.Printf("  FileName:  %s\n", item.FileName)

		log.Println("----------------------------------------")
	}

	log.Println("========================================")
	fmt.Println()
}