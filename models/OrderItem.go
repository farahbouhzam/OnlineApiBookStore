package main

type OrderItem struct { 
    Book     Book `json:"book"` 
    Quantity int  `json:"quantity"` 
}