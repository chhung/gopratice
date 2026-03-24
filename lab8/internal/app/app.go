package app

import (
	"context"
	"flag"
	"fmt"
	"io"

	"lab8/internal/config"
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

	if _, err := fmt.Fprintf(stdout, "ActiveMQ consumer started, config=%s\n", *configPath); err != nil {
		return err
	}

	return consumer.Run(ctx)
}
