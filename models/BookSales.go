package models



type BookSales struct { 
	
	BookID int
    Quantity int  `json:"quantity_sold"`
} 
