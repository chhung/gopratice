package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/go-stomp/stomp/v3"
	"go.uber.org/zap"

	"lab8/internal/config"
	"lab8/internal/logging"
)

type Consumer struct {
	cfg    config.Config
	stdout io.Writer
	logger *zap.Logger
	conn   *stomp.Conn
	subs   []*namedSubscription
}

type namedSubscription struct {
	label        string
	destination  string
	subscription *stomp.Subscription
}

func NewConsumer(cfg config.Config, stdout io.Writer) (*Consumer, error) {
	if stdout == nil {
		return nil, fmt.Errorf("stdout writer is required")
	}

	logger := logging.NewJSONLogger(stdout, "consumer")

	return &Consumer{cfg: cfg, stdout: stdout, logger: logger}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := c.connect(); err != nil {
		return err
	}
	defer c.close()

	c.logger.Info("connected to ActiveMQ", zap.String("address", c.cfg.Address()))

	if err := c.subscribe(); err != nil {
		return err
	}

	c.logger.Info("waiting for messages")

	errCh := make(chan error, len(c.subs))
	var wg sync.WaitGroup
	for _, sub := range c.subs {
		wg.Add(1)
		go func(current *namedSubscription) {
			defer wg.Done()
			c.consume(runCtx, current, errCh)
		}(sub)
	}

	var runErr error
	select {
	case <-runCtx.Done():
	case runErr = <-errCh:
	}

	cancel()
	wg.Wait()

	if runErr != nil {
		return runErr
	}

	return nil
}

func (c *Consumer) connect() error {
	options := []func(*stomp.Conn) error{
		stomp.ConnOpt.HeartBeat(5*time.Second, 5*time.Second),
		stomp.ConnOpt.Host(c.hostHeader()),
	}

	if c.cfg.Broker.Username != "" || c.cfg.Broker.Password != "" {
		options = append(options, stomp.ConnOpt.Login(c.cfg.Broker.Username, c.cfg.Broker.Password))
	}

	conn, err := stomp.Dial("tcp", c.cfg.Address(), options...)
	if err != nil {
		return fmt.Errorf("connect to ActiveMQ %s: %w", c.cfg.Address(), err)
	}

	c.conn = conn
	return nil
}

func (c *Consumer) subscribe() error {
	destinations := []struct {
		label       string
		destination string
	}{
		{label: "queue", destination: c.cfg.QueueDestination()},
		{label: "topic", destination: c.cfg.TopicDestination()},
	}

	for _, item := range destinations {
		if item.destination == "" {
			continue
		}

		subscription, err := c.conn.Subscribe(item.destination, stomp.AckAuto,
			stomp.SubscribeOpt.Header("activemq.prefetchSize", "5000"),
		)
		if err != nil {
			return fmt.Errorf("subscribe %s %s: %w", item.label, item.destination, err)
		}

		c.subs = append(c.subs, &namedSubscription{
			label:        item.label,
			destination:  item.destination,
			subscription: subscription,
		})

		c.logger.Info("subscribed to destination",
			zap.String("label", item.label),
			zap.String("destination", item.destination),
		)
	}

	return nil
}

func (c *Consumer) consume(ctx context.Context, sub *namedSubscription, errCh chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-sub.subscription.C:
			if !ok {
				select {
				case errCh <- fmt.Errorf("subscription closed: %s", sub.destination):
				default:
				}
				return
			}

			if msg == nil {
				continue
			}

			if msg.Err != nil {
				select {
				case errCh <- fmt.Errorf("receive %s message: %w", sub.destination, msg.Err):
				default:
				}
				return
			}

			body := strings.TrimRight(string(msg.Body), "\r\n")
			c.logger.Info("received message",
				zap.String("label", sub.label),
				zap.String("destination", sub.destination),
				zap.String("body", body),
			)
		}
	}
}

func (c *Consumer) close() error {
	var closeErr error

	for _, sub := range c.subs {
		if sub == nil || sub.subscription == nil {
			continue
		}

		if err := sub.subscription.Unsubscribe(); err != nil && !errors.Is(err, stomp.ErrCompletedSubscription) {
			closeErr = errors.Join(closeErr, fmt.Errorf("unsubscribe %s: %w", sub.destination, err))
		}
	}

	if c.conn != nil {
		if err := c.conn.Disconnect(); err != nil && !errors.Is(err, stomp.ErrAlreadyClosed) {
			closeErr = errors.Join(closeErr, fmt.Errorf("disconnect: %w", err))
		}
	}

	return closeErr
}

func (c *Consumer) hostHeader() string {
	if c.cfg.Broker.HostName != "" {
		return c.cfg.Broker.HostName
	}

	return c.cfg.Broker.Host
}
