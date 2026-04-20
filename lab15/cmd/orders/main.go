package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Order struct {
	VendorID   string    `bson:"vendor_id"`
	CustomerID string    `bson:"customer_id"`
	OrderDate  time.Time `bson:"order_date"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 48*time.Hour)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Fatal("connect error:", err)
	}
	defer client.Disconnect(ctx)

	vendorIDs := generateVendorIDs(1, 10000)
	customerIDs := generateCustomerIDs(100000)
	batchSize := 5000

	totalVendors := len(vendorIDs)
	for i, vid := range vendorIDs {
		err = batchInsertOrders(ctx, client, vid, customerIDs, batchSize)
		if err != nil {
			log.Fatalf("batch insert error for vendor %s: %v", vid, err)
		}
		fmt.Printf("vendor %d/%d done (%s)\n", i+1, totalVendors, vid)
	}
	fmt.Printf("insert completed, total vendors: %d, total orders: %d\n", totalVendors, totalVendors*len(customerIDs))
}

func batchInsertOrders(ctx context.Context, client *mongo.Client, vendorID string, customerIDs []string, batchSize int) error {
	collection := client.Database("customerAggregationDB").Collection("orders")

	total := len(customerIDs)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batch := customerIDs[i:end]

		now := time.Now()
		docs := make([]interface{}, len(batch))
		for j, cid := range batch {
			docs[j] = Order{
				VendorID:   vendorID,
				CustomerID: cid,
				OrderDate:  now,
			}
		}

		_, err := collection.InsertMany(ctx, docs)
		if err != nil {
			return fmt.Errorf("insert failed at batch %d-%d: %w", i, end-1, err)
		}
	}
	return nil
}

func generateVendorIDs(startFrom, endInclusive int) []string {
	ids := make([]string, endInclusive-startFrom+1)
	for i := 0; i < len(ids); i++ {
		ids[i] = fmt.Sprintf("V%014d", startFrom+i)
	}
	return ids
}

func generateCustomerIDs(n int) []string {
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = fmt.Sprintf("C%014d", i+1)
	}
	return ids
}
