package smbget

import (
	"io"
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	var debug = os.Getenv("DEBUG") == "true"
	var handler *slog.JSONHandler

	if debug {
		opts := &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}

		handler = slog.NewJSONHandler(os.Stderr, opts)

	} else {
		handler = slog.NewJSONHandler(io.Discard, nil)
	}

	logger = slog.New(handler)
}
