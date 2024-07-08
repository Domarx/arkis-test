package main

import (
	"arkis_test/queue"
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	address     = flag.String("address", "amqp://localhost:5672", "rabbitmq address")
	inputQueue  = flag.String("input", "", "name of the queue to pulish messages to")
	outputQueue = flag.String("output", "", "name of the queue to read messages from")
)

func main() {
	flag.Parse()

	if *inputQueue == "" && *outputQueue == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	wg := sync.WaitGroup{}

	if *inputQueue != "" {
		publisher, err := queue.New(*address, *inputQueue)
		if err != nil {
			logrus.WithError(err).Panic("Cannot create input queue")
		}

		scanner := bufio.NewScanner(os.Stdin)

		wg.Add(1)
		go func() {
			defer wg.Done()

			for scanner.Scan() {
				message := scanner.Bytes()
				if err = publisher.Publish(context.Background(), message); err != nil {
					logrus.WithError(err).Error("failed to publish message")
					return
				}
			}

			if err = scanner.Err(); err != nil {
				logrus.WithError(err).Error("error scanning input")
			}
		}()

	}

	if *outputQueue != "" {
		consumer, err := queue.New(*address, *outputQueue)
		if err != nil {
			logrus.WithError(err).Panic("Cannot create input queue")
		}

		messages, err := consumer.Consume(context.Background())
		if err != nil {
			logrus.WithError(err).Panic("failed to consume messages")
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			for message := range messages {
				fmt.Fprintln(os.Stdout, "[recv] ", string(message.Body))
			}
		}()
	}

	wg.Wait()
}
