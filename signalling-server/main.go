package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	envErr := godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Now that logging is configured, we can log any errors from loading the .env file
	if envErr != nil {
		slog.Info("No .env file found, relying on environment variables")
	}

	// Start the HTTP server.
	if err := Serve(8089); err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
		return
	}
}
