package repository

import (
	"context"
	"testing"

	"lab2/internal/model"
)

func TestInMemoryMessageRepositoryCreate(t *testing.T) {
	repo := NewInMemoryMessageRepository([]model.Message{{ID: 3, Text: "existing"}})

	created, err := repo.Create(context.Background(), "new message")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if created.ID != 4 {
		t.Fatalf("Create().ID = %d, want 4", created.ID)
	}

	listed, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(listed) != 2 {
		t.Fatalf("List() len = %d, want 2", len(listed))
	}
}
