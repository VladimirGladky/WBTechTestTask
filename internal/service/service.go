package service

import (
	"WBTechTestTask/internal/models"
	"WBTechTestTask/internal/repository"
	"WBTechTestTask/pkg/logger"
	"WBTechTestTask/pkg/suberrors"
	"context"
	"fmt"
	"sync"
)

type OrderServiceInterface interface {
	GetOrder(id string) (*models.Order, error)
	Create(order *models.Order) (string, error)
	InitCache(ctx context.Context) error
}

type OrderService struct {
	repo  repository.OrderRepositoryInterface
	ctx   context.Context
	cache map[string]*models.Order
	mu    *sync.RWMutex
}

func (o OrderService) InitCache(ctx context.Context) error {
	logger.GetLoggerFromCtx(ctx).Info("init cache")
	orders, err := o.repo.GetAllOrders()
	if err != nil {
		return err
	}

	o.mu.Lock()
	defer o.mu.Unlock()
	o.cache = make(map[string]*models.Order, len(orders))
	for _, order := range orders {
		o.cache[order.OrderUid] = order
	}
	return nil
}

func (o OrderService) GetOrder(id string) (*models.Order, error) {
	if id == "" {
		return nil, fmt.Errorf("order id is empty")
	}
	o.mu.RLock()
	order, ok := o.cache[id]
	o.mu.RUnlock()
	if ok {
		return order, nil
	}
	order, err := o.repo.GetOrder(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, suberrors.ErrIdOrderNotFound
	}
	o.mu.Lock()
	o.cache[id] = order
	o.mu.Unlock()
	return order, nil
}

func (o OrderService) Create(order *models.Order) (string, error) {
	if order == nil {
		return "", fmt.Errorf("order is empty")
	}
	if order.OrderUid == "" {
		return "", fmt.Errorf("order uid is empty")
	}
	id, err := o.repo.CreateOrder(order)
	if err != nil {
		return "", err
	}
	return id, nil
}

func NewOrderService(repo repository.OrderRepositoryInterface, ctx context.Context) *OrderService {
	return &OrderService{
		repo:  repo,
		ctx:   ctx,
		cache: make(map[string]*models.Order),
		mu:    &sync.RWMutex{},
	}
}
