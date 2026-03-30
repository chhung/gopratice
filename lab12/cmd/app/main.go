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
	//address     = "10.40.138.16:61613"
	address     = "127.0.0.1:61613"
	destination = "/topic/momo_coupon"
	username    = "admin"
	password    = "admin"
)

func main() {
	// 這個跟activemq無法，主要是用來捕捉Ctrl+C信號，讓程式能夠優雅地關閉連線和訂閱。
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// 當接收到中斷信號時，stop()會被調用，從而取消ctx，讓run函數知道應該結束運行。
	// 這個等到main()結束時才會調用，確保在run函數中能夠正確地捕捉到ctx.Done()信號。
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
		case <-ctx.Done(): // 當接收到中斷信號時，ctx.Done()會被觸發，這裡我們打印一條消息並返回nil，讓程式優雅地結束。
			fmt.Println("shutting down")
			return nil
		case msg, ok := <-sub.C: // 從訂閱的通道中接收消息。如果通道被關閉，ok會是false，這意味著訂閱已經意外地關閉了，我們應該返回一個錯誤。
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
