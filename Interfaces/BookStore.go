package interfaces


import (
	"online_bookStore/models"
)
type BookStore interface { 
 CreateBook(book models.Book) (models.Book, error) 
 GetBook(id int) (models.Book, error) 
 UpdateBook(id int, book models.Book) (models.Book, error) 
 DeleteBook(id int) error 
 
} 