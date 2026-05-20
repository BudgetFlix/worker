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
)

func main() {

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	cfg := config.Load()

	connection := mustRabbitConnection(cfg)
	defer connection.Close()

	consumer := rabbitmq.NewConsumer(
		connection,
	)

	msgs, err := consumer.Consume(
		"video.upload.queue",
	)

	if err != nil {
		log.Fatal(err)
	}

	go rabbitmq.Loop(
		ctx,
		msgs,
		handler.Media,
	)

	waitShutdown(ctx)
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