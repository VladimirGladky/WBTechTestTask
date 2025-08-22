package service

import "WBTechTestTask/internal/models"

type OrderServiceInterface interface {
	ProcessOrder(order models.Order) error
}
