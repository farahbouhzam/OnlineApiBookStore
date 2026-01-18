package concreteimplemetations

import (
	"context"
	"database/sql"
	"encoding/json"
	"online_bookStore/models"
	
)

type MySQLBookStore struct {
	db *sql.DB
	
}

// Constructore

func NewMySQLBookStore(db *sql.DB) *MySQLBookStore {
	return &MySQLBookStore{
		db: db,
	}
}

func (s *MySQLBookStore) CreateBook(ctx context.Context,book models.Book) (models.Book, error) {

	

	genresJson, err := json.Marshal(book.Genres)

	if err != nil {
		return book, err
	}

	query := `
		INSERT INTO books (title, genres, published_at, price, stock, author_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(
		ctx,
		query,
		book.Title,
		string(genresJson),
		book.PublishedAt,
		book.Price,
		book.Stock,
		book.Author.ID,
	)

	if err != nil {
		return book, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return book, err
	}

	book.ID = int(id)
	return book, nil
}

func (s *MySQLBookStore) GetBook(ctx context.Context, id int) (models.Book, error){

    query := `
		SELECT
			b.id, b.title, b.genres, b.published_at, b.price, b.stock,
			a.id, a.first_name, a.last_name, a.bio
		FROM books b
		JOIN authors a ON b.author_id = a.id
		WHERE b.id = ?
	`

	var book models.Book
	var genresJson string


	err := s.db.QueryRowContext(ctx,query, id).Scan(
		&book.ID,                 // book ID
		&book.Title,              // book title
		&genresJson,              // genres as JSON string
		&book.PublishedAt,        // publication date
		&book.Price,              // price
		&book.Stock,              // stock
		&book.Author.ID,          // author ID
		&book.Author.FirstName,   // author first name
		&book.Author.LastName,    // author last name
		&book.Author.Bio,         // author bio
	)
    
	if err != nil {
		return book, err
	}


	err = json.Unmarshal([]byte(genresJson),&book.Genres)

	if err != nil {
		return book, err
	}


	return book, nil

}

func (s *MySQLBookStore) UpdateBook(ctx context.Context, id int, book models.Book) (models.Book, error){

	genresJSON, err := json.Marshal(book.Genres)
	if err != nil {
		return book, err
	}



	query := `
		UPDATE books
		SET title = ?, genres = ?, published_at = ?, price = ?, stock = ?, author_id = ?
		WHERE id = ?
	`


	result, err := s.db.ExecContext(
		ctx,
		query,
		book.Title,
		string(genresJSON),
		book.PublishedAt,
		book.Price,
		book.Stock,
		book.Author.ID,
		id,
	)

	if err != nil {
		return book, err
	}

	rowsAffected, err := result.RowsAffected()
    if err !=nil {
		return book, err
	}


	if rowsAffected==0{
		return book, sql.ErrNoRows
	}


	book.ID=id

	return book, nil

}


func (s *MySQLBookStore) DeleteBook(ctx context.Context,id int) error {
	 query := `
		DELETE FROM books
		WHERE id = ?
	`

	result, err:= s.db.ExecContext(ctx,query,id)

	if err !=nil{
		return err
	} 

	rowAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}


	if rowAffected == 0 {
		return sql.ErrNoRows
	}


	return nil
}




func (s *MySQLBookStore) SearchBooks(ctx context.Context,searchCriteria models.SearchCriteria) ([]models.Book, error) {

	var books []models.Book

	query := `
		SELECT b.id, b.title, b.published_at, b.price, b.stock, b.author_id
		FROM books b
		WHERE 1 = 1
	`

	var args []interface{}

	if searchCriteria.Title != "" {
		query += " AND b.title LIKE ?"
		args = append(args, "%"+searchCriteria.Title+"%")
	}

	if searchCriteria.AuthorID != 0 {
		query += " AND b.author_id = ?"
		args = append(args, searchCriteria.AuthorID)
	}

	if searchCriteria.Genre != "" {
		query += " AND b.genres LIKE ?"
		args = append(args, "%"+searchCriteria.Genre+"%")
	}

	if searchCriteria.MinPrice != 0 {
		query += " AND b.price >= ?"
		args = append(args, searchCriteria.MinPrice)
	}

	if searchCriteria.MaxPrice != 0 {
		query += " AND b.price <= ?"
		args = append(args, searchCriteria.MaxPrice)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var book models.Book

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.PublishedAt,
			&book.Price,
			&book.Stock,
			&book.Author.ID,
		)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}
