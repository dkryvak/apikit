package main

import (
	"apikit/internal/app"
	"apikit/internal/config"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.NewApp(config.Version()).Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
