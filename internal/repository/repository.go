package repository

import "WBTechTestTask/internal/models"

type OrderRepositoryInterface interface {
	GetOrder(order *models.Order) error
	CreateOrder(order *models.Order) error
}
