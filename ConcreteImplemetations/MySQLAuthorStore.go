package concreteimplemetations

import (
	"database/sql"

	"online_bookStore/models"
	
)

type MySQLAuthorStore struct {
	db *sql.DB
	
}

// Constructore

func NewMySQLAuthorStore(db *sql.DB) *MySQLAuthorStore {
	return &MySQLAuthorStore{
		db: db,
	}
}

func (s *MySQLAuthorStore) CreateAuthor(author models.Author) (models.Author, error) {

	

	
	query := `
		INSERT INTO authors (first_name, last_name, bio)
		VALUES (?, ?, ?)
	`

	result, err := s.db.Exec(
		query,
		author.FirstName,
		author.LastName,
		author.Bio,
	)

	if err != nil {
		return author, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return author, err
	}

	author.ID = int(id)
	return author, nil
}

func (s *MySQLAuthorStore) GetAuthor(id int) (models.Author, error){

    query := `
		SELECT id, first_name, last_name, bio
		FROM authors
		WHERE id = ?
	`

	var author models.Author
	


	err := s.db.QueryRow(query, id).Scan(
		&author.ID,
		&author.FirstName,
		&author.LastName,
		&author.Bio,
	)
    
	if err != nil {
		return author, err
	}

	return author, nil

}

func (s *MySQLAuthorStore) UpdateAuthor(id int, author models.Author) (models.Author, error){

	


	query := `
		UPDATE authors
		SET first_name = ?, last_name = ?, bio = ?
		WHERE id = ?
	`


	result, err := s.db.Exec(
		query,
		author.FirstName,
		author.LastName,
		author.Bio,
		id,
	)

	if err != nil {
		return author, err
	}

	rowsAffected, err := result.RowsAffected()
    if err !=nil {
		return author, err
	}


	if rowsAffected==0{
		return author, sql.ErrNoRows
	}


	author.ID=id

	return author, nil

}


func (s *MySQLAuthorStore) DeleteAuthor(id int) error {
	 query := `
		DELETE FROM authors
		WHERE id = ?
	`

	result, err:= s.db.Exec(query,id)

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


