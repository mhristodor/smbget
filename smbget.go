package smbget

import (
	"io"
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	var debug = os.Getenv("DEBUG") == "true"
	var handler *slog.TextHandler

	if debug {
		opts := &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}

		handler = slog.NewTextHandler(os.Stderr, opts)

	} else {
		handler = slog.NewTextHandler(io.Discard, nil)
	}

	logger = slog.New(handler)
}
