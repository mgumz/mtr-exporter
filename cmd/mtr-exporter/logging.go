package main

import (
	"fmt"
	"log/slog"
	"os"
)

func logLevel(level string) (*slog.LevelVar, error) {
	ll := &slog.LevelVar{}
	switch level {
	case "debug":
		ll.Set(slog.LevelDebug)
	case "info":
		ll.Set(slog.LevelInfo)
	case "warn":
		ll.Set(slog.LevelWarn)
	case "error":
		ll.Set(slog.LevelError)
	default:
		return nil, fmt.Errorf("unsupported -log-level: %q", level)
	}
	return ll, nil
}

// removes time attr from logging.
// copied from log/slog/internal/slogtest
func removeTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}

func setupLogging(level string, useTimestamps bool) {

	logLevel, err := logLevel(level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	logOptions := slog.HandlerOptions{
		Level: logLevel,
	}

	if !useTimestamps {
		logOptions.ReplaceAttr = removeTime
	}

	logHandler := slog.NewTextHandler(os.Stdout, &logOptions)
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
}
