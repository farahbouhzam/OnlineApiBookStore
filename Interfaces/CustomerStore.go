package interfaces


import (
	"online_bookStore/models"
)
type CustomerStore interface {
	CreateCustomer(customer models.Customer) (models.Customer, error)
	GetCustomer(id int) (models.Customer, error)
	UpdateCustomer(id int, customer models.Customer) (models.Customer, error)
	DeleteCustomer(id int) error
}