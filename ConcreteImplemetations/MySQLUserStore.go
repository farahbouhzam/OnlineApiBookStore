package concreteimplemetations

import (
	"context"
	"database/sql"
	"online_bookStore/models"
)

type MySQLUserStore struct {
	db *sql.DB
}

func NewMySQLUserStore(db *sql.DB) *MySQLUserStore {
	return &MySQLUserStore{db: db}
}


func (s *MySQLUserStore) GetByEmail(
	ctx context.Context,
	email string,
) (models.User, error) {

	query := `
		SELECT id, email, password, role
		FROM users
		WHERE email = ?
	`

	var user models.User
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *MySQLUserStore) GetByID(
	ctx context.Context,
	id int,
) (models.User, error) {

	query := `
		SELECT id, email, password, role
		FROM users
		WHERE id = ?
	`

	var user models.User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *MySQLUserStore) CreateUser(
	ctx context.Context,
	user models.User,
) (models.User, error) {

	query := `
		INSERT INTO users (email, password, role)
		VALUES (?, ?, ?)
	`

	result, err := s.db.ExecContext(
		ctx,
		query,
		user.Email,
		user.Password,
		user.Role,
	)
	if err != nil {
		return user, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return user, err
	}

	user.ID = int(id)
	return user, nil
}


