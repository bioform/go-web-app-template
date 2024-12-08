package logging

import (
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/bioform/go-web-app-template/pkg/env"
	"github.com/lmittmann/tint"
)

var once sync.Once

func InitLogger() {
	once.Do(func() {

		opts, err := handlerOptions()
		logger := slog.New(handler(opts))

		if err != nil {
			logger.Error("Error initializing logger", slog.Any("error", err))
		}

		slog.SetDefault(logger)

		slog.Info("Logger initialized", slog.Any("level", opts.Level))
	})
}

func handlerOptions() (*slog.HandlerOptions, error) {
	level, err := logLevel()

	wd, wdErr := os.Getwd()

	return &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				if wdErr == nil {
					truncateSourcePath(wd, a)
				}
			}
			return a
		},
	}, err
}

func handler(opts *slog.HandlerOptions) slog.Handler {
	if env.IsProduction() {
		return slog.NewJSONHandler(os.Stdout, opts)
	}

	if colorize() {
		return tint.NewHandler(os.Stdout, &tint.Options{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: opts.ReplaceAttr,
			TimeFormat:  time.TimeOnly,
		})
	}

	return slog.NewTextHandler(os.Stdout, opts)
}

func truncateSourcePath(wd string, a slog.Attr) {
	s := a.Value.Any().(*slog.Source)
	file, err := filepath.Rel(wd, s.File)
	if err != nil {
		return
	}
	s.File = file
}

func colorize() bool {
	colorize := env.Get("LOG_COLORIZE", "true")
	return colorize == "true" || colorize == "1"
}
