package concreteimplemetations

import (
	"context"
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

func (s *MySQLOrderStore) CreateOrder(ctx context.Context, order models.Order) (models.Order, error) {


    tx, err := s.db.BeginTx(ctx, nil)
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

    result, err := tx.ExecContext(
        ctx,
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

    itemQuery := `
        INSERT INTO order_items (order_id, book_id, quantity)
        VALUES (?, ?, ?)
    `

    for _, item := range order.Items {
        _, err = tx.ExecContext(
            ctx,
            itemQuery,
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

func (s *MySQLOrderStore) GetOrder(ctx context.Context, id int) (models.Order, error) {
	query := `
		SELECT 
			o.id, o.total_price, o.created_at, o.status,
			c.id, c.name, c.email
		FROM orders o
		JOIN customers c ON o.customer_id = c.id
		WHERE o.id = ?
	`

	var order models.Order
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.TotalPrice,
		&order.CreatedAt,
		&order.Status,
		&order.Customer.ID,
		&order.Customer.Name,
		&order.Customer.Email,
	)
	if err != nil {
		return order, err
	}

	itemsQuery := `
		SELECT 
			oi.id, oi.quantity,
			b.id, b.title, b.genres, b.published_at, b.price, b.stock
		FROM order_items oi
		JOIN books b ON oi.book_id = b.id
		WHERE oi.order_id = ?
	`

	rows, err := s.db.QueryContext(ctx, itemsQuery, order.ID)
	if err != nil {
		return order, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		var genresJSON string

		err := rows.Scan(
			&item.ID,
			&item.Quantity,
			&item.Book.ID,
			&item.Book.Title,
			&genresJSON,
			&item.Book.PublishedAt,
			&item.Book.Price,
			&item.Book.Stock,
		)
		if err != nil {
			return order, err
		}

		_ = json.Unmarshal([]byte(genresJSON), &item.Book.Genres)
		order.Items = append(order.Items, item)
	}

	return order, nil
}


func (s *MySQLOrderStore) UpdateOrderStatus(ctx context.Context,id int, status string) (models.Order, error) {
	var order models.Order

	query := `
	  UPDATE orders
	  SET status = ?		
	  WHERE id = ?
	`

	result, err := s.db.ExecContext(ctx,query, status, id)
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

func (s *MySQLOrderStore) DeleteOrder(ctx context.Context,id int) error {
	query := `
	    DELETE FROM orders
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

func (s *MySQLOrderStore) GetOrderByDateRange(
	ctx context.Context,
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

	rows, err := s.db.QueryContext(ctx,query, from, to)
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


func (s *MySQLOrderStore) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	query := `
	   SELECT o.id, o.total_price, o.created_at, o.status,
	   c.id, c.name,c.email

	   From orders o
	   JOIN customers c ON o.customer_id = c.id

	`

	rows , err := s.db.QueryContext(ctx,query)

	if err != nil {
		return nil , err
	}

	var orders []models.Order

	for rows.Next(){
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