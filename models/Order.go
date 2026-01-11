package models


import (
	"time"
)

type Order struct { 
    ID         int          
    Customer   Customer    
    Items      []OrderItem 
    TotalPrice float64     
    CreatedAt  time.Time    
    Status     string       
}