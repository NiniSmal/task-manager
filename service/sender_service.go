package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type SenderService struct {
	kafka *kafka.Writer
}

func NewSenderService(w *kafka.Writer) *SenderService {
	return &SenderService{
		kafka: w,
	}
}

type Email struct {
	Text    string `json:"text"`
	To      string `json:"to"`
	Subject string `json:"subject"`
}

func (s *SenderService) SendEmail(ctx context.Context, email Email) error {
	msg, err := json.Marshal(&email)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	err = s.kafka.WriteMessages(ctx, kafka.Message{Value: msg})
	if err != nil {
		return fmt.Errorf("write messages: %w", err)
	}

	return nil
}
