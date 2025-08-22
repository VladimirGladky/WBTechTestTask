package service

import "WBTechTestTask/internal/models"

type OrderServiceInterface interface {
	GetOrder(id string) (models.Order, error)
	Create(order models.Order) (string, error)
}
