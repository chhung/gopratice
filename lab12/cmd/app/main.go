package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/go-stomp/stomp/v3"
)

const (
	address = "10.40.138.16:61613"
	//address          = "127.0.0.1:61613"
	destination = "/topic/momo_coupon"
	username    = "admin"
	password    = "admin"
	clientID    = "momo-coupon-go-consumer"
	//subscriptionName = "momo-coupon-go-subscription"
	subscriptionName = "momo-coupon-go-durable"
	workerCount      = 20
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
		stomp.ConnOpt.Header("client-id", clientID),
	)
	if err != nil {
		return fmt.Errorf("connect to ActiveMQ %s: %w", address, err)
	}
	defer conn.Disconnect()

	fmt.Printf("connected to %s\n", address)

	sub, err := conn.Subscribe(destination, stomp.AckAuto,
		stomp.SubscribeOpt.Header("activemq.prefetchSize", "5000"),
		stomp.SubscribeOpt.Header("activemq.subscriptionName", subscriptionName),
	)
	if err != nil {
		return fmt.Errorf("subscribe %s: %w", destination, err)
	}
	defer sub.Unsubscribe()

	fmt.Printf("listening on %s — press Ctrl+C to exit\n", destination)

	msgCh := make(chan *stomp.Message, workerCount*2)

	var wg sync.WaitGroup
	for range workerCount {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range msgCh {
				process(msg)
			}
		}()
	}
	defer func() {
		close(msgCh)
		wg.Wait()
	}()

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
			select {
			case msgCh <- msg:
			case <-ctx.Done():
				fmt.Println("shutting down")
				return nil
			}
		}
	}
}

func process(msg *stomp.Message) {
	consumerReceivedAt := time.Now()
	brokerReceivedAtMs, _ := strconv.ParseInt(msg.Header.Get("timestamp"), 10, 64)
	brokerReceivedAt := time.UnixMilli(brokerReceivedAtMs)

	const isoFmt = "2006-01-02T15:04:05.000-07:00"
	fmt.Printf("activemq received at: %s | consumer received at: %s | latency: %d ms | body: %s\n",
		brokerReceivedAt.Format(isoFmt), consumerReceivedAt.Format(isoFmt),
		consumerReceivedAt.Sub(brokerReceivedAt).Milliseconds(), string(msg.Body))
}
