package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ExportSnapshotID struct {
	VendorID      string `bson:"vendor_id"`
	ConditionType string `bson:"condition_type"`
	WindowField   string `bson:"window_field"`
	Threshold     int    `bson:"threshold"`
}

type ExportSnapshot struct {
	ID            ExportSnapshotID `bson:"_id"`
	VendorID      string           `bson:"vendor_id"`
	ConditionType string           `bson:"condition_type"`
	WindowField   string           `bson:"window_field"`
	Threshold     int              `bson:"threshold"`
	CustomerCount int              `bson:"customer_count"`
	CustomerIDs   string           `bson:"customer_ids"`
	GeneratedAt   time.Time        `bson:"generated_at"`
	ExpiresAt     time.Time        `bson:"expires_at"`
	Status        string           `bson:"status"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://mongoadmin:secret@127.0.0.1:27017"))
	if err != nil {
		log.Fatal("connect error:", err)
	}
	defer client.Disconnect(ctx)

	customerIDs := generateCodes("C", 100000)
	vendorIDs := generateVendorIDs(327574, 520000)
	batchSize := 500

	total := len(vendorIDs)
	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batch := vendorIDs[i:end]

		err = batchInsertExportSnapshots(ctx, client, batch, "order", "last_180", 4, customerIDs)
		if err != nil {
			log.Fatalf("batch insert error at %d-%d: %v", i, end-1, err)
		}
		fmt.Printf("inserted %d / %d\n", end, total)
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Printf("insert completed, total: %d\n", total)
}

func batchInsertExportSnapshots(ctx context.Context, client *mongo.Client, vendorIDs []string, conditionType, windowField string, threshold int, customerIDs string) error {
	collection := client.Database("customerAggregationDB").Collection("export_snapshots")

	now := time.Now()
	expires := now.Add(24 * time.Hour)

	docs := make([]interface{}, len(vendorIDs))
	for i, vid := range vendorIDs {
		docs[i] = ExportSnapshot{
			ID: ExportSnapshotID{
				VendorID:      vid,
				ConditionType: conditionType,
				WindowField:   windowField,
				Threshold:     threshold,
			},
			VendorID:      vid,
			ConditionType: conditionType,
			WindowField:   windowField,
			Threshold:     threshold,
			CustomerCount: 100000,
			CustomerIDs:   customerIDs,
			GeneratedAt:   now,
			ExpiresAt:     expires,
			Status:        "completed",
		}
	}

	_, err := collection.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("batch insert failed: %w", err)
	}
	return nil
}

func generateCodes(prefix string, n int) string {
	codes := make([]string, n)
	for i := 0; i < n; i++ {
		codes[i] = fmt.Sprintf("%s%014d", prefix, i+1)
	}
	return strings.Join(codes, ",")
}

func generateVendorIDs(startFrom, endInclusive int) []string {
	ids := make([]string, endInclusive-startFrom+1)
	for i := 0; i < len(ids); i++ {
		ids[i] = fmt.Sprintf("V%014d", startFrom+i)
	}
	return ids
}
