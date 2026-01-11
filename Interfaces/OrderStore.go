package interfaces

import (
	"online_bookStore/models"
	"time"
	"context"
)

type OrderStore interface {
	CreateOrder(ctx context.Context,order models.Order) (models.Order, error)
	GetOrder(ctx context.Context,id int) (models.Order, error)
	UpdateOrderStatus(ctx context.Context,id int, status string) (models.Order, error)
	DeleteOrder(ctx context.Context,id int) error
	GetOrderByDateRange(ctx context.Context,from time.Time, to time.Time) ([]models.Order, error)
	GetAllOrders(ctx context.Context) ([]models.Order, error)
	
}
