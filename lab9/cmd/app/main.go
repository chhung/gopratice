package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func ConnectDB() *mongo.Client {
	// 1. 讀取 .env 檔
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoURI := os.Getenv("MONGODB_URI")

	// 2. 設定連線選項
	clientOptions := options.Client().ApplyURI(mongoURI)

	// 3. 建立連線 (Context 設定超時避免卡死)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 4. 檢查連線是否真的成功 (Ping)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	fmt.Println("Successfully connected to MongoDB!")
	return client
}

func main() {
	client := ConnectDB()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("disconnect mongodb: %v", err)
		}
	}()

	fmt.Println("MongoDB client is ready")
}
