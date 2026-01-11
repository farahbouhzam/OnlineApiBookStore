package concreteimplemetations

import (
	"database/sql"
	"encoding/json"
	"time"

	"online_bookStore/models"
)

type MySQLOrderStore struct {
	db *sql.DB
}

func NewMySQLOrderStore(db *sql.DB) *MySQLOrderStore {
	return &MySQLOrderStore{
		db: db,
	}
}

func (s *MySQLOrderStore) CreateOrder(order models.Order) (models.Order, error) {

	tx, err := s.db.Begin()

	if err != nil {
		return order, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	orderQuery := `
        INSERT INTO orders (customer_id, total_price, status)
        VALUES (?, ?, ?)
    `

	result, err := tx.Exec(
		orderQuery,
		order.Customer.ID,
		order.TotalPrice,
		order.Status,
	)

	if err != nil {
		tx.Rollback()
		return order, err
	}

	orderID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return order, err
	}

	order.ID = int(orderID)

	query := `
	INSERT INTO order_items (order_id, book_id, quantity)
	VALUES (?, ?, ?)
	`

	for _, item := range order.Items {
		_, err = tx.Exec(
			query,
			order.ID,
			item.Book.ID,
			item.Quantity,
		)

		if err != nil {
			tx.Rollback()
			return order, err
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()

		return order, err
	}

	return order, nil

}

func (s *MySQLOrderStore) GetOrder(id int) (models.Order, error) {

	query := `
	SELECT
	  o.id, o.customer_id, o.total_price, o.created_at, o.status, 
      c.id, c.name, c.email
	FROM orders o
	JOIN customers c ON o.customer_id = c.id
	WHERE o.id=?
	`

	var order models.Order

	err := s.db.QueryRow(query, id).Scan(
		&order.ID,
		&order.Customer.ID,
		&order.TotalPrice,
		&order.CreatedAt,
		&order.Status,
	)

	if err != nil {
		return order, err
	}

	// getting order items

	queryItem := `
	  SELECT 
	    oi.id, oi.quantity,
		b.id, b.title, b.genres, b.published_at, b.price, b.stock
		FROM order_items oi
		JOIN books b ON oi.book_id = b.id
		WHERE oi.order_id = ?

	`

	rows, err := s.db.Query(queryItem, order.ID)

	if err != nil {
		return order, err
	}

	defer rows.Close()
	for rows.Next() {
		var item models.OrderItem
		var jsonGenre string
		err := rows.Scan(
			&item.ID,
			&item.Quantity,
			&item.Book.ID,
			&item.Book.Title,
			&item.Book.Title,
			&jsonGenre,
			&item.Book.PublishedAt,
			&item.Book.Price,
			&item.Book.Stock,
		)

		if err != nil {
			return order, err
		}

		err = json.Unmarshal([]byte(jsonGenre), &item.Book.Genres)
		if err != nil {
			return order, err
		}

		order.Items = append(order.Items, item)
	}

	return order, nil

}

func (s *MySQLOrderStore) UpdateOrderStatus(id int, status string) (models.Order, error) {
	var order models.Order

	query := `
	  UPDATE orders
	  SET status = ?		
	  WHERE id = ?
	`

	result, err := s.db.Exec(query, status, id)
	if err != nil {
		return order, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return order, err
	}

	if rowsAffected == 0 {
		return order, sql.ErrNoRows
	}

	order.ID = id
	order.Status = status

	return order, nil

}

func (s *MySQLOrderStore) DeleteOrder(id int) error {
	query := `
	    DELETE FROM orders
		WHERE id = ?
	 `

	result, err := s.db.Exec(query, id)

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

func (s *MySQLOrderStore) GetOrderByDateRange(
	from time.Time,
	to time.Time,
) ([]models.Order, error) {

	query := `
		SELECT 
			o.id, o.total_price, o.created_at, o.status,
			c.id, c.name, c.email
		FROM orders o
		JOIN customers c ON o.customer_id = c.id
		WHERE o.created_at BETWEEN ? AND ?
		ORDER BY o.created_at ASC
	`

	rows, err := s.db.Query(query, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order

		err := rows.Scan(
			&order.ID,
			&order.TotalPrice,
			&order.CreatedAt,
			&order.Status,
			&order.Customer.ID,
			&order.Customer.Name,
			&order.Customer.Email,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}


