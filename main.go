package main

import (
	"arkis_test/database"
	"arkis_test/processor"
	"arkis_test/queue"
	"context"
	"os"
	"os/signal"
	"sync"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	wg := sync.WaitGroup{}

	for _, config := range []struct {
		inputQueue  string
		outputQueue string
	}{
		{inputQueue: "input-A", outputQueue: "output-A"},
		{inputQueue: "input-B", outputQueue: "output-B"},
	} {
		inputQueue, err := queue.New(os.Getenv("RABBITMQ_URL"), config.inputQueue)
		if err != nil {
			log.WithError(err).Panic("Cannot create input queue")
		}

		outputQueue, err := queue.New(os.Getenv("RABBITMQ_URL"), config.outputQueue)
		if err != nil {
			log.WithError(err).Panic("Cannot create output queue")
		}

		wg.Add(1)
		go func() {
			processor.New(inputQueue, outputQueue, database.D{}).Run(ctx)
			wg.Done()
		}()
	}

	log.Info("Application is ready to run")
	wg.Wait()
}
