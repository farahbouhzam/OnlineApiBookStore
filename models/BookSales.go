package models



type BookSales struct {
	BookID   int `json:"book_id"`
	Quantity int `json:"quantity_sold"`
}
