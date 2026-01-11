package interfaces

import (
	"online_bookStore/models"
	"time"
)

type OrderStore interface {
	CreateOrder(order models.Order) (models.Order, error)
	GetOrder(id int) (models.Order, error)
	UpdateOrderStatus(id int, status string) (models.Order, error)
	DeleteOrder(id int) error
	GetOrderByDateRange(from time.Time, to time.Time) ([]models.Order, error)
}
