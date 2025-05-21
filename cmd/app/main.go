package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"task_1/internal/filepresenter"
	"task_1/internal/fileproducer"
	"task_1/internal/service"
)

func setLogger(logLevel string) *slog.Logger {
	levelVar := new(slog.LevelVar)
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: levelVar},
		),
	)
	switch logLevel {
	case "debug":
		levelVar.Set(slog.LevelDebug)
	case "info":
		levelVar.Set(slog.LevelInfo)
	case "warn":
		levelVar.Set(slog.LevelWarn)
	case "error":
		levelVar.Set(slog.LevelError)
	default:
		levelVar.Set(slog.LevelInfo)
	}
	slog.SetDefault(logger)
	return logger
}

func main() {
	var pathProduce, pathPresent string

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var logger *slog.Logger
	var logLevel string
	cmd := &cli.Command{
		Name:    "url-masking",
		Usage:   "Application that can mask URLs to an asterisks with logging",
		Version: "v0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Value:       "info",
				Usage:       "log level",
				Destination: &logLevel,
			},
			&cli.StringFlag{
				Name:        "path",
				Usage:       "path to the input file",
				Required:    true,
				Destination: &pathProduce,
			},
		},
		Action: func(ctx context.Context, command *cli.Command) error {
			logger = setLogger(logLevel)
			logger.Info("Logger initialized", "level", logLevel, "--path", pathProduce)
			return nil
		},
	}

	if err := cmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	go func() {
		for {
			select {
			case <-exit:
				cancel()
				logger.Error("Program was interrupted...")
				return
			}
		}
	}()

	producer, err := fileproducer.NewFileProducer(pathProduce)
	if err != nil {
		logger.Error(fmt.Sprintf("Producer error: %v", err))
		return
	}
	presenter := filepresenter.NewFilePresenter(pathPresent)
	newService := service.NewService(producer, presenter)
	err = newService.Run(ctx, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("Service cannot to start. Error: %v", err))
		return
	}

	logger.Info("Program finished", "path", pathProduce)
}
