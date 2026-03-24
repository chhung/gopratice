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

	"lab8/internal/config"
)

type Consumer struct {
	cfg    config.Config
	stdout io.Writer
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

	return &Consumer{cfg: cfg, stdout: stdout}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := c.connect(); err != nil {
		return err
	}
	defer c.close()

	if _, err := fmt.Fprintf(c.stdout, "Connected to ActiveMQ at %s\n", c.cfg.Address()); err != nil {
		return err
	}

	if err := c.subscribe(); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(c.stdout, "Waiting for messages. Press Ctrl+C to stop."); err != nil {
		return err
	}

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

		subscription, err := c.conn.Subscribe(item.destination, stomp.AckAuto)
		if err != nil {
			return fmt.Errorf("subscribe %s %s: %w", item.label, item.destination, err)
		}

		c.subs = append(c.subs, &namedSubscription{
			label:        item.label,
			destination:  item.destination,
			subscription: subscription,
		})

		if _, err := fmt.Fprintf(c.stdout, "Subscribed %s: %s\n", item.label, item.destination); err != nil {
			return err
		}
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
			if _, err := fmt.Fprintf(c.stdout, "[%s] %s => %s\n", sub.label, sub.destination, body); err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
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
