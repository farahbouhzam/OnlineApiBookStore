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

```bash
go run main.go

Expected logs:

Starting Online Bookstore API
Database connected
Server running on :8081


If port 8081 is already in use, change it in main.go.


