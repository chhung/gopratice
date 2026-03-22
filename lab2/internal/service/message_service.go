package service

import (
	"context"
	"errors"
	"strings"

	"lab2/internal/model"
	"lab2/internal/repository"
)

var ErrEmptyMessage = errors.New("message text is required")
var ErrMessageTooLong = errors.New("message text exceeds 280 characters")

const maxMessageLength = 280

type MessageService struct {
	repository repository.MessageRepository
}

func NewMessageService(messageRepository repository.MessageRepository) *MessageService {
	return &MessageService{repository: messageRepository}
}

func (s *MessageService) List(ctx context.Context) ([]model.Message, error) {
	return s.repository.List(ctx)
}

func (s *MessageService) Create(ctx context.Context, text string) (model.Message, error) {
	normalizedText := strings.TrimSpace(text)
	if normalizedText == "" {
		return model.Message{}, ErrEmptyMessage
	}

	if len([]rune(normalizedText)) > maxMessageLength {
		return model.Message{}, ErrMessageTooLong
	}

	return s.repository.Create(ctx, normalizedText)
}
