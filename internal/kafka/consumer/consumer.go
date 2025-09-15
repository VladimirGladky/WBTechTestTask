package kafka

import (
	"WBTechTestTask/internal/models"
	"WBTechTestTask/internal/service"
	"WBTechTestTask/pkg/logger"
	"context"
	"encoding/json"
	"errors"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type Consumer struct {
	reader   *kafka.Reader
	ctx      context.Context
	orderSrv service.OrderServiceInterface
	cancel   context.CancelFunc
}

func NewConsumer(brokers []string, topic string, orderSrv service.OrderServiceInterface, ctx context.Context) *Consumer {
	ctx, cancel := context.WithCancel(ctx)
	return &Consumer{
		ctx:      ctx,
		orderSrv: orderSrv,
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
		cancel: cancel,
	}
}

func (c *Consumer) ConsumeMessages() error {
	for {
		select {
		case <-c.ctx.Done():
			return c.Close()
		default:
			select {
			case <-c.ctx.Done():
				return c.Close()
			default:
				msg, err := c.reader.ReadMessage(c.ctx)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						return c.Close()
					}
					logger.GetLoggerFromCtx(c.ctx).Error("Error reading message: %v", zap.Error(err))
					return err
				}
				if err := c.ProcessMessage(msg); err != nil {
					logger.GetLoggerFromCtx(c.ctx).Error("Error processing message: %v", zap.Error(err))
					continue
				}
				logger.GetLoggerFromCtx(c.ctx).Info("Message processed", zap.String("message", string(msg.Value)))
			}
		}
	}
}

func (c *Consumer) ProcessMessage(msg kafka.Message) error {
	var order *models.Order
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		return err
	}
	if err := order.Validate(); err != nil {
		return err
	}
	if _, err := c.orderSrv.Create(order); err != nil {
		return err
	}
	return nil
}

func (c *Consumer) Close() error {
	c.cancel()
	return c.reader.Close()
}
