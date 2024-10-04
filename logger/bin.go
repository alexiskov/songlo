package logger

import (
	"io"
	"log/slog"
)

var Log *slog.Logger

func New(w io.Writer) {
	Log = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(Log)
}
