package main

import (
	"fmt"
	"online_bookStore/DataBase"
)



func main() {
	db := database.NewMySQLDBFromEnv()
	defer db.Close()

	fmt.Println("MySQL connected ")
}
