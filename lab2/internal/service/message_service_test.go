package service

import (
	"context"
	"errors"
	"testing"

	"lab2/internal/model"
	"lab2/internal/repository"
)

func TestMessageServiceCreate(t *testing.T) {
	service := NewMessageService(repository.NewInMemoryMessageRepository([]model.Message{
		{ID: 1, Text: "hello"},
		{ID: 2, Text: "world"},
	}))

	created, err := service.Create(context.Background(), "new message")
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if created.ID != 3 {
		t.Fatalf("Create().ID = %d, want 3", created.ID)
	}

	if created.Text != "new message" {
		t.Fatalf("Create().Text = %q, want %q", created.Text, "new message")
	}

	messages, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(messages) != 3 {
		t.Fatalf("List() len = %d, want 3", len(messages))
	}
}

func TestMessageServiceCreateRejectsEmptyText(t *testing.T) {
	service := NewMessageService(repository.NewInMemoryMessageRepository(nil))

	_, err := service.Create(context.Background(), "   ")
	if !errors.Is(err, ErrEmptyMessage) {
		t.Fatalf("Create() error = %v, want %v", err, ErrEmptyMessage)
	}
}
