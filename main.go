package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"worker/internal/config"
	"worker/internal/handler"
	"worker/internal/rabbitmq"
	"worker/internal/storage"
)

func main() {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg := config.Load()

	recoverProcessingJobs(cfg)

	connection := mustRabbitConnection(cfg)
	defer connection.Close()

	consumer := rabbitmq.NewConsumer(
		connection,
	)

	msgs, err := consumer.Consume(
		cfg.UploadQueue,
	)

	if err != nil {
		log.Fatal(err)
	}

	producer := rabbitmq.NewProducer(
		connection,
	)

	go rabbitmq.Loop(
		ctx,
		msgs,
		handler.Media,
		producer,
		&cfg,
	)

	waitShutdown(ctx)
}

func recoverProcessingJobs(
	cfg config.Config,
) {
	if err := storage.RecoverProcessingJobs(cfg); err != nil {
		log.Printf("failed to recover processing jobs: %v", err)
	}
}

func mustRabbitConnection(
	cfg config.Config,
) *rabbitmq.Connection {

	connection, err := rabbitmq.NewConnection(
		cfg.RabbitMQURL(),
	)

	if err != nil {
		log.Fatal(err)
	}

	return connection
}

func waitShutdown(
	ctx context.Context,
) {

	<-ctx.Done()

	log.Println("shutdown signal received")

	done := make(chan struct{})

	go func() {
		handler.PipelineWG.Wait()
		close(done)
	}()

	select {

	case <-done:
		log.Println("all pipelines finished")

	case <-time.After(3 * time.Minute): // three minute for the testing
		log.Println("shutdown timeout reached")
	}
}
