package repository

import (
	"context"
	"sync"

	"lab2/internal/model"
)

type MessageRepository interface {
	List(ctx context.Context) ([]model.Message, error)
	Create(ctx context.Context, text string) (model.Message, error)
}

type InMemoryMessageRepository struct {
	mu       sync.RWMutex
	messages []model.Message
	nextID   int
}

func NewInMemoryMessageRepository(seed []model.Message) *InMemoryMessageRepository {
	messages := make([]model.Message, len(seed))
	copy(messages, seed)

	nextID := 1
	for _, message := range messages {
		if message.ID >= nextID {
			nextID = message.ID + 1
		}
	}

	return &InMemoryMessageRepository{
		messages: messages,
		nextID:   nextID,
	}
}

func (r *InMemoryMessageRepository) List(_ context.Context) ([]model.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.Message, len(r.messages))
	copy(result, r.messages)
	return result, nil
}

func (r *InMemoryMessageRepository) Create(_ context.Context, text string) (model.Message, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	message := model.Message{
		ID:   r.nextID,
		Text: text,
	}

	r.messages = append(r.messages, message)
	r.nextID++

	return message, nil
}
