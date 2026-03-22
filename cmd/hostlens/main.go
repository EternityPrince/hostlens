package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"hostlens/internal/app"
)

func run() error {
	startupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	application, err := app.New(startupCtx, os.Args[1:])
	if err != nil {
		return fmt.Errorf("error when init app: %w", err)
	}

	defer application.Close()
	return application.Run(context.Background())
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
