package interfaces


import (
	"online_bookStore/models"
	"context"
)
type CustomerStore interface {
	CreateCustomer(ctx context.Context,customer models.Customer) (models.Customer, error)
	GetCustomer(ctx context.Context,id int) (models.Customer, error)
	UpdateCustomer(ctx context.Context,id int, customer models.Customer) (models.Customer, error)
	DeleteCustomer(ctx context.Context,id int) error
	GetAllCustomers(ctx context.Context) ([]models.Customer, error)
}