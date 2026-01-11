package interfaces


import (
	"online_bookStore/models"
	"context"
)
type BookStore interface { 
 CreateBook(ctx context.Context,book models.Book) (models.Book, error) 
 GetBook(ctx context.Context,id int) (models.Book, error) 
 UpdateBook(ctx context.Context,id int, book models.Book) (models.Book, error) 
 DeleteBook(ctx context.Context,id int) error 
 SearchBooks(ctx context.Context,searchCriteria models.SearchCriteria)([]models.Book, error)
} 