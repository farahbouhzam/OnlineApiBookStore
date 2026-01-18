package interfaces

import (
	"context"
	"online_bookStore/models"
)

type UserStore interface {
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetByID(ctx context.Context, id int) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
}

