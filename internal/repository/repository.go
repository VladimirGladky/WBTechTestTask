package repository

import "WBTechTestTask/internal/models"

type OrderRepositoryInterface interface {
	Save(order *models.Order) error
}
