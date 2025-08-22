package repository

import "WBTechTestTask/internal/models"

type OrderRepositoryInterface interface {
	GetOrder(id string) (*models.Order, error)
	CreateOrder(order *models.Order) (string, error)
}
