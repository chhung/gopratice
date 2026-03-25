package app

import (
	"context"
	"flag"
	"io"

	"go.uber.org/zap"

	"lab8/internal/config"
	"lab8/internal/logging"
	"lab8/internal/service"
)

const defaultConfigPath = "config.json"

func Run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	flagSet := flag.NewFlagSet("lab8-consumer", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	configPath := flagSet.String("config", defaultConfigPath, "path to the ActiveMQ config file")
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		return err
	}

	consumer, err := service.NewConsumer(cfg, stdout)
	if err != nil {
		return err
	}

	logger := logging.NewJSONLogger(stdout, "app")
	logger.Info("ActiveMQ consumer started", zap.String("config", *configPath))

	return consumer.Run(ctx)
}
