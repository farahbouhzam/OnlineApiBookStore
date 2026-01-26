package models


type OrderItem struct {
	ID       int  `json:"id"`
	Book     Book `json:"book"`
	Quantity int  `json:"quantity"`
}
