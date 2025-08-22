package service

import (
	"WBTechTestTask/internal/models"
	"WBTechTestTask/internal/repository"
	"context"
	"fmt"
)

type OrderServiceInterface interface {
	GetOrder(id string) (*models.Order, error)
	Create(order *models.Order) (string, error)
}

type OrderService struct {
	repo repository.OrderRepositoryInterface
	ctx  context.Context
}

func (o OrderService) GetOrder(id string) (*models.Order, error) {
	if id == "" {
		return nil, fmt.Errorf("order id is empty")
	}
	return o.repo.GetOrder(id)
}

func (o OrderService) Create(order *models.Order) (string, error) {
	if order == nil {
		return "", fmt.Errorf("order is empty")
	}
	if order.OrderId == "" {
		return "", fmt.Errorf("order id is empty")
	}
	return o.repo.CreateOrder(order)
}

func NewOrderService(repo repository.OrderRepositoryInterface, ctx context.Context) OrderServiceInterface {
	return &OrderService{
		repo: repo,
		ctx:  ctx,
	}
}
