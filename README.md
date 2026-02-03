# Online Bookstore API (Go + MySQL)

This project is a RESTful Online Bookstore API built with Go and MySQL. It follows a clean, modular structure with separate layers for handlers, services, models, and data access.

## Project Structure
```
ConcreteImplemetations/   MySQL data access layer
DataBase/                 DB connection + SQL schema/seed
Handlers/                 HTTP handlers
Interfaces/               Interfaces
models/                   Domain models
reports/                  Generated JSON reports
services/                 Business services (sales reports)
main.go                   Application entry point
go.mod                    Go module file
```

## Features Implemented
- Authors CRUD
- Books CRUD
- Books search by title, genre, author, price range
- Customers CRUD with addresses
- Orders CRUD with multiple items
- Transaction-safe order creation
- Daily sales report generation (JSON)
- Background job with graceful shutdown
- Context usage with timeouts
- Basic logging of key events

## Requirements
- Go 1.25.5+
- MySQL 8.x (or MariaDB compatible)

## Environment Variables
Set these before running:
```
DB_USER=root
DB_PASSWORD=your_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=online_bookstore
```

## Database Setup (Windows)
Create database:
```powershell
mysql -u root -p -h localhost -P 3306
```
```sql
CREATE DATABASE online_bookstore;
```

Create tables and insert sample data:
```powershell
mysql -u root -p -h localhost -P 3306 online_bookstore < DataBase\schema.sql
mysql -u root -p -h localhost -P 3306 online_bookstore < DataBase\data.sql
```

## Database Setup (macOS)
Install MySQL (Homebrew):
```bash
brew install mysql
brew services start mysql
```

Create database:
```bash
mysql -u root -p -h localhost -P 3306
```
```sql
CREATE DATABASE online_bookstore;
```

Create tables and insert sample data:
```bash
mysql -u root -p -h localhost -P 3306 online_bookstore < DataBase/schema.sql
mysql -u root -p -h localhost -P 3306 online_bookstore < DataBase/data.sql
```

## Run the API
```powershell
go run main.go
```
Expected logs:
```
Starting Online Bookstore API
Database connected
Server running on :8081
```

If port 8081 is in use, change it in `main.go`.

## API Testing (PowerShell Examples)
Author create:
```powershell
Invoke-RestMethod `
  -Uri "http://localhost:8081/authors" `
  -Method POST `
  -Headers @{ "Content-Type"="application/json" } `
  -Body '{
    "first_name": "George",
    "last_name": "Orwell",
    "bio": "English novelist"
  }'
```

Book create:
```powershell
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
```

Customer create:
```powershell
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
```

Order create:
```powershell
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
```

## Reports
- A sales report is generated every 24 hours by a background job.
- Files are saved under `reports/` as `sales_report_YYYY-MM-DD.json`.
- Reports API:
  - `GET /reports` list report files
  - `GET /reports/{YYYY-MM-DD}` fetch a report by date

Note: The first report is generated after 24 hours (the job uses a 24h ticker).

## Common Endpoints
- `GET /authors`, `POST /authors`, `GET /authors/{id}`, `PUT /authors/{id}`, `DELETE /authors/{id}`
- `GET /books`, `POST /books`, `GET /books/{id}`, `PUT /books/{id}`, `DELETE /books/{id}`
- `GET /customers`, `POST /customers`, `GET /customers/{id}`, `PUT /customers/{id}`, `DELETE /customers/{id}`
- `GET /orders`, `POST /orders`, `GET /orders/{id}`, `PUT /orders/{id}`, `DELETE /orders/{id}`
- `GET /reports`, `GET /reports/{date}`
