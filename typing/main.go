package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-stomp/stomp/v3"
)

const (
	address     = "10.40.138.16:61613"
	destination = "/topic/momo_coupon"
	username    = "admin"
	password    = "admin"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	conn, err := stomp.Dial("tcp", address,
		stomp.ConnOpt.Login(username, password),
		stomp.ConnOpt.HeartBeat(5*time.Second, 5*time.Second),
	)
	if err != nil {
		return fmt.Errorf("connect to ActiveMQ %s: %w", address, err)
	}
	defer conn.Disconnect()

	fmt.Printf("connected to %s\n", address)

	sub, err := conn.Subscribe(destination, stomp.AckAuto)
	if err != nil {
		return fmt.Errorf("subscribe %s: %w", destination, err)
	}
	defer sub.Unsubscribe()

	fmt.Printf("listening on %s — press Ctrl+C to exit\n", destination)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("shutting down")
			return nil
		case msg, ok := <-sub.C:
			if !ok {
				return fmt.Errorf("subscription closed unexpectedly")
			}
			if msg.Err != nil {
				return fmt.Errorf("receive message: %w", msg.Err)
			}

			consumerReceivedAt := time.Now().UnixMilli()
			brokerReceivedAt, _ := strconv.ParseInt(msg.Header.Get("timestamp"), 10, 64)

			fmt.Printf("activemq received at: %d ms | consumer received at: %d ms | latency: %d ms | body: %s\n",
				brokerReceivedAt, consumerReceivedAt, consumerReceivedAt-brokerReceivedAt, string(msg.Body))
		}
	}
}
