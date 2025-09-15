package service

import (
	"WBTechTestTask/internal/models"
	"WBTechTestTask/internal/repository"
	"WBTechTestTask/pkg/logger"
	"WBTechTestTask/pkg/suberrors"
	"context"
	"fmt"
	"go.uber.org/zap"
	"sync"
	"time"
)

type OrderServiceInterface interface {
	GetOrder(id string) (*models.Order, error)
	Create(order *models.Order) (string, error)
	InitCache(ctx context.Context) error
}

type OrderService struct {
	repo  repository.OrderRepositoryInterface
	ctx   context.Context
	cache map[string]*CacheEntry
	mu    *sync.RWMutex
	ttl   time.Duration
}

type CacheEntry struct {
	order *models.Order
	ts    time.Time
}

func NewOrderService(repo repository.OrderRepositoryInterface, ctx context.Context) *OrderService {
	return &OrderService{
		repo:  repo,
		ctx:   ctx,
		cache: make(map[string]*CacheEntry),
		mu:    &sync.RWMutex{},
		ttl:   time.Second * 30,
	}
}

func (o *OrderService) InitCache(ctx context.Context) error {
	orders, err := o.repo.GetAllOrders()
	if err != nil {
		return err
	}

	o.mu.Lock()
	defer o.mu.Unlock()
	o.cache = make(map[string]*CacheEntry, len(orders))
	for _, order := range orders {
		o.cache[order.OrderUid] = &CacheEntry{
			order: order,
			ts:    time.Now(),
		}
	}
	logger.GetLoggerFromCtx(ctx).Info("cache init completed", zap.Any("orders", len(orders)))
	return nil
}

func (o *OrderService) GetOrder(id string) (*models.Order, error) {
	if id == "" {
		return nil, fmt.Errorf("order id is empty")
	}
	o.mu.RLock()
	ce, ok := o.cache[id]
	o.mu.RUnlock()
	if ok {
		return ce.order, nil
	}
	order, err := o.repo.GetOrder(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, suberrors.ErrIdOrderNotFound
	}
	o.mu.Lock()
	o.cache[id] = &CacheEntry{
		order: order,
		ts:    time.Now(),
	}
	o.mu.Unlock()
	return order, nil
}

func (o *OrderService) Create(order *models.Order) (string, error) {
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

func (o *OrderService) StartEviction(interval string) {
	t, err := time.ParseDuration(interval)
	if err != nil {
		logger.GetLoggerFromCtx(o.ctx).Error("error parsing interval", zap.Error(err))
	}
	ticker := time.NewTicker(t)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			o.evict()
		case <-o.ctx.Done():
			return
		}
	}
}

func (o *OrderService) evict() {
	o.mu.Lock()
	defer o.mu.Unlock()
	now := time.Now()
	for id, entry := range o.cache {
		if now.Sub(entry.ts) > o.ttl {
			logger.GetLoggerFromCtx(o.ctx).Info("evicting order", zap.String("order_id", id))
			delete(o.cache, id)
		}
	}
}
