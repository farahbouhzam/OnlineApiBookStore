package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"online_bookStore/database"
	"online_bookStore/concreteimplemetations"
	"online_bookStore/handlers"
	"online_bookStore/services"
)

func main() {
	log.Println("Starting Online Bookstore API")

	// Root context (for shutdown + background jobs)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ---- DATABASE ----
	db := database.NewMySQLDBFromEnv()
	defer db.Close()
	log.Println("Database connected")

	// ---- STORES ----
	authorStore := concreteimplemetations.NewMySQLAuthorStore(db)
	bookStore := concreteimplemetations.NewMySQLBookStore(db)
	customerStore := concreteimplemetations.NewMySQLCustomerStore(db)
	orderStore := concreteimplemetations.NewMySQLOrderStore(db)
	userStore := concreteimplemetations.NewMySQLUserStore(db)

	_ = userStore // used later for JWT auth

	// ---- SERVICES ----
	salesReportService := services.NewSalesReportService(orderStore)

	// ---- BACKGROUND JOBS ----
	services.StartSalesReportJob(ctx, salesReportService)

	// ---- HANDLERS ----
	authorHandler := handlers.NewAuthorHandler(authorStore)
	bookHandler := handlers.NewBookHandler(bookStore)
	customerHandler := handlers.NewCustomerHandler(customerStore)
	orderHandler := handlers.NewOrderHandler(orderStore)
	reportHandler := handlers.NewReportHandler()

	// ---- ROUTES ----
	mux := http.NewServeMux()

	mux.HandleFunc("/authors", authorHandler.AuthorsHandler)
	mux.HandleFunc("/authors/", authorHandler.AuthorsByIDHandler)

	mux.HandleFunc("/books", bookHandler.BooksHandler)
	mux.HandleFunc("/books/", bookHandler.BookByIDHandler)

	mux.HandleFunc("/customers", customerHandler.CustomersHandler)
	mux.HandleFunc("/customers/", customerHandler.CustomersByIDHandler)

	mux.HandleFunc("/orders", orderHandler.OrdersHandler)
	mux.HandleFunc("/orders/", orderHandler.OrdersByIDHandler)

	mux.HandleFunc("/reports", reportHandler.GetReports)
	mux.HandleFunc("/reports/", reportHandler.GetReportByDate)

	// ---- SERVER ----
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server
	go func() {
		log.Println("Server running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// ---- GRACEFUL SHUTDOWN ----
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Forced shutdown: %v", err)
	}

	cancel() // stop background jobs
	log.Println("Server stopped cleanly")
}
