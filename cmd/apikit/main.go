package main

import (
	"apikit/internal/app"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// version is injected at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := app.NewApp(version).Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
