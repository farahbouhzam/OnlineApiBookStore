package database


import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLDBFromEnv() *sql.DB {
	user := os.("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	if user == "" || password == "" || host == "" || port == "" || dbname == "" {
		log.Fatal("database environment variables are not set")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user, password, host, port, dbname,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	return db
}
