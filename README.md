# 📚 Online Bookstore API – Go

This project is a **RESTful Online Bookstore API** built with **Go** and **MySQL**, following clean architecture principles.  
It supports book, author, customer, and order management, advanced search, background processing, logging, and sales reporting.

---

## 🎯 Project Objectives

- Build a clean and modular REST API in Go
- Apply good backend practices (context, logging, transactions)
- Manage concurrency and background tasks
- Generate periodic business reports
- Expose data via JSON APIs

---

## ✅ Features Implemented

### 📖 Book Management
- Create, read, update, delete books
- Search books by:
  - title
  - genre
  - author
  - price range

### ✍️ Author Management
- Full CRUD operations for authors

### 👤 Customer Management
- Full CRUD operations
- Address management

### 🛒 Order Management
- Create orders with multiple items
- Update order status
- Delete orders
- Fetch orders
- Transaction-safe order creation

### 📊 Sales Reports
- Automatic **daily sales report**
- Aggregates:
  - total revenue
  - number of orders
  - top-selling books
- Saved as **JSON files**
- Accessible via API

### 🧵 Concurrency & Reliability
- `context.Context` with timeouts
- Graceful cancellation using `ctx.Done()`
- Background goroutine for periodic tasks

### 🪵 Logging
- API errors
- Database failures
- Business events:
  - book created
  - order placed
  - report generated

---

## 🏗 Project Structure

├── concreteimplemetations/ # MySQL data access layer
├── database/ # DB connection logic
├── handlers/ # HTTP handlers
├── interfaces/ # Interfaces
├── models/ # Domain models
├── services/ # Sales report service
├── reports/ # Generated JSON reports
├── main.go # Application entry point
├── go.mod
└── README.md


---

## 🛠 Technologies Used

- Go
- MySQL
- database/sql
- go-sql-driver/mysql
- REST architecture
- JSON
- Goroutines

---

## ⚙️ Environment Configuration

Set the following environment variables:

DB_USER=root
DB_PASSWORD=your_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=online_bookstore


---

## ▶️ How to Run the Project

``` bash
go run main.go

Expected logs:

Starting Online Bookstore API
Database connected
Server running on :8081


If port 8081 is already in use, change it in main.go.


API Testing (PowerShell Examples)
➕ Create an Author
Invoke-RestMethod `
  -Uri "http://localhost:8081/authors" `
  -Method POST `
  -Headers @{ "Content-Type"="application/json" } `
  -Body '{
    "first_name": "George",
    "last_name": "Orwell",
    "bio": "English novelist"
  }'

➕ Create a Book
Invoke-RestMethod `
  -Uri "http://localhost:8081/books" `
  -Method POST `
  -Headers @{ "Content-Type"="application/json" } `
  -Body '{
    "title": "1984",
    "genres": ["Dystopian", "Political"],
    "published_at": "1949-06-08T00:00:00Z",
    "price": 19.99,
    "stock": 10,
    "author": { "id": 1 }
  }'

🔍 Search Books
Invoke-RestMethod "http://localhost:8081/books?title=1984&min_price=10"

➕ Create a Customer
Invoke-RestMethod `
  -Uri "http://localhost:8081/customers" `
  -Method POST `
  -Headers @{ "Content-Type"="application/json" } `
  -Body '{
    "name": "Alice",
    "email": "alice@example.com",
    "address": {
      "street": "Main St",
      "city": "Paris",
      "state": "IDF",
      "postal_code": "75000",
      "country": "France"
    }
  }'

➕ Create an Order
Invoke-RestMethod `
  -Uri "http://localhost:8081/orders" `
  -Method POST `
  -Headers @{ "Content-Type"="application/json" } `
  -Body '{
    "customer": { "id": 1 },
    "status": "pending",
    "total_price": 39.98,
    "items": [
      {
        "book": { "id": 1 },
        "quantity": 2
      }
    ]
  }'

🔄 Update Order Status
Invoke-RestMethod `
  -Uri "http://localhost:8081/orders/1" `
  -Method PUT `
  -Headers @{ "Content-Type"="application/json" } `
  -Body '{ "status": "completed" }'

📊 Sales Reports
🕒 Background Job

Runs every 24 hours

Generates a sales report

Saves it as a JSON file in /reports

Example:

reports/sales_report_2026-01-26.json

📡 Access Reports via API
List reports
Invoke-RestMethod "http://localhost:8081/reports"

Get report by date
Invoke-RestMethod "http://localhost:8081/reports/2026-01-26"

🧵 Concurrency & Context

Every request uses context.WithTimeout

Database queries use QueryContext / ExecContext

Background goroutine stops gracefully with ctx.Done()

Transactions used for order creation

🪵 Logging

Logged events include:

API errors

Database failures

Business actions:

BOOK CREATED

ORDER CREATED

CUSTOMER UPDATED

SALES REPORT SAVED