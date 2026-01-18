package interfaces

import (
	"context"
	"online_bookStore/models"
)

type AuthorStore interface {
	CreateAuthor(ctx context.Context, author models.Author) (models.Author, error)
	GetAuthor(ctx context.Context, id int) (models.Author, error)
	UpdateAuthor(ctx context.Context, id int, author models.Author) (models.Author, error)
	DeleteAuthor(ctx context.Context, id int) error
	GetAllAuthors(ctx context.Context) ([]models.Author, error)
	
}

