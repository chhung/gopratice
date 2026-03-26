package app

import (
	"context"
	"flag"
	"fmt"
	"io"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	"lab9/internal/config"
)

const defaultConfigPath = "configs/config.yaml"

func Run(ctx context.Context, args []string, stdout, stderr io.Writer) error {
	flagSet := flag.NewFlagSet("lab9", flag.ContinueOnError)
	flagSet.SetOutput(stderr)

	configPath := flagSet.String("config", defaultConfigPath, "path to the config file")
	if err := flagSet.Parse(args); err != nil {
		return err
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		return err
	}

	connCtx, cancel := context.WithTimeout(ctx, cfg.MongoTimeout())
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		return fmt.Errorf("connect to mongodb: %w", err)
	}

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			fmt.Fprintf(stderr, "disconnect mongodb: %v\n", err)
		}
	}()

	if err := client.Ping(connCtx, readpref.Primary()); err != nil {
		return fmt.Errorf("ping mongodb: %w", err)
	}

	fmt.Fprintln(stdout, "Successfully connected to MongoDB!")
	fmt.Fprintf(stdout, "Database: %s\n", cfg.MongoDB.Database)

	<-ctx.Done()
	return nil
}
