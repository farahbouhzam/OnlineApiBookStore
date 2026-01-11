package concreteimplemetations

import (
	"database/sql"

	"online_bookStore/models"
	"context"
	
)

type MySQLCustomerStore struct {
	db *sql.DB
}

// Constructor
func NewMySQLCustomerStore(db *sql.DB) *MySQLCustomerStore {
	return &MySQLCustomerStore{
		db: db,
	}
}


func (s *MySQLCustomerStore) CreateCustomer(ctx context.Context, customer models.Customer) (models.Customer, error) {

	// 1. Insert address
	addressQuery := `
		INSERT INTO addresses (street, city, state, postal_code, country)
		VALUES (?, ?, ?, ?, ?)
	`

	addressResult, err := s.db.ExecContext(
		ctx,
		addressQuery,
		customer.Address.Street,
		customer.Address.City,
		customer.Address.State,
		customer.Address.PostalCode,
		customer.Address.Country,
	)
	if err != nil {
		return customer, err
	}

	addressID, err := addressResult.LastInsertId()
	if err != nil {
		return customer, err
	}

	// 2. Insert customer
	customerQuery := `
		INSERT INTO customers (name, email, address_id)
		VALUES (?, ?, ?)
	`

	result, err := s.db.Exec(
		customerQuery,
		customer.Name,
		customer.Email,
		addressID,
	)
	if err != nil {
		return customer, err
	}

	customerID, err := result.LastInsertId()
	if err != nil {
		return customer, err
	}

	customer.ID = int(customerID)
	customer.Address.ID = int(addressID)

	return customer, nil
}

func (s *MySQLCustomerStore) GetCustomer(ctx context.Context,id int) (models.Customer, error) {

	query := `
		SELECT
			c.id, c.name, c.email, c.created_at,
			a.id, a.street, a.city, a.state, a.postal_code, a.country
		FROM customers c
		JOIN addresses a ON c.address_id = a.id
		WHERE c.id = ?
	`

	var customer models.Customer

	err := s.db.QueryRowContext(ctx,query, id).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.CreatedAt,
		&customer.Address.ID,
		&customer.Address.Street,
		&customer.Address.City,
		&customer.Address.State,
		&customer.Address.PostalCode,
		&customer.Address.Country,
	)
	if err != nil {
		return customer, err
	}

	return customer, nil
}


func (s *MySQLCustomerStore) UpdateCustomer(ctx context.Context,id int, customer models.Customer) (models.Customer, error) {

	// 1. Update address
	addressQuery := `
		UPDATE addresses
		SET street = ?, city = ?, state = ?, postal_code = ?, country = ?
		WHERE id = ?
	`

	_, err := s.db.ExecContext(
		ctx,
		addressQuery,
		customer.Address.Street,
		customer.Address.City,
		customer.Address.State,
		customer.Address.PostalCode,
		customer.Address.Country,
		customer.Address.ID,
	)
	if err != nil {
		return customer, err
	}

	// 2. Update customer
	customerQuery := `
		UPDATE customers
		SET name = ?, email = ?
		WHERE id = ?
	`

	result, err := s.db.Exec(
		customerQuery,
		customer.Name,
		customer.Email,
		id,
	)
	if err != nil {
		return customer, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return customer, err
	}

	if rowsAffected == 0 {
		return customer, sql.ErrNoRows
	}

	customer.ID = id
	return customer, nil
}


func (s *MySQLCustomerStore) DeleteCustomer(ctx context.Context,id int) error {

	query := `
		DELETE FROM customers
		WHERE id = ?
	`

	result, err := s.db.ExecContext(ctx,query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}





