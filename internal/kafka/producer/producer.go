package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
	topic  string
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.RoundRobin{},
		},
		topic: topic,
	}
}

func (p *Producer) SendOrderMessage(ctx context.Context, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	message := kafka.Message{
		Value: jsonValue,
		Time:  time.Now(),
	}

	err = p.writer.WriteMessages(ctx, message)
	log.Println("Message sent to topic", p.topic)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
