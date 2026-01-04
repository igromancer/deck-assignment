package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/igromancer/deck-assignment/internal/queue"
)

func main() {
	fmt.Println("starting job processor")

	jobReceiver, err := queue.NewJobReceiver()
	if err != nil {
		panic(fmt.Errorf("failed to create the job processor: %s", err.Error()))
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	err = jobReceiver.Run(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to start the job processor:  %s", err.Error()))
	}
}
