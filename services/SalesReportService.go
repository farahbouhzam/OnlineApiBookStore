package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"online_bookStore/Interfaces"
	"online_bookStore/models"
)

type SalesReportService struct {
	orderStore interfaces.OrderStore
}

// Constructor
func NewSalesReportService(orderStore interfaces.OrderStore) *SalesReportService {
	return &SalesReportService{
		orderStore: orderStore,
	}
}

// Save report to JSON file
func SaveReportToJson(report models.SalesReport) error {
	// ensure directory exists
	if err := os.MkdirAll("reports", 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf(
		"sales_report_%s.json",
		report.Timestamp.Format("2006-01-02"),
	)

	path := filepath.Join("reports", filename)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(report)
}

// Generate report logic
func (s *SalesReportService) GenerateSalesReport(
	ctx context.Context,
	from time.Time,
	to time.Time,
) (models.SalesReport, error) {

	orders, err := s.orderStore.GetOrderByDateRange(ctx, from, to)
	if err != nil {
		return models.SalesReport{}, err
	}

	totalRevenue := 0.0
	totalOrders := len(orders)

	bookCounter := make(map[int]int) // bookID → quantity sold

	for _, order := range orders {
		totalRevenue += order.TotalPrice

		for _, item := range order.Items {
			bookCounter[item.Book.ID] += item.Quantity
		}
	}

	var sales []models.BookSales
	for bookID, qty := range bookCounter {
		sales = append(sales, models.BookSales{
			BookID:   bookID,
			Quantity: qty,
		})
	}

	sort.Slice(sales, func(i, j int) bool {
		return sales[i].Quantity > sales[j].Quantity
	})

	limit := 3
	if len(sales) < limit {
		limit = len(sales)
	}

	topSellingBooks := sales[:limit]

	return models.SalesReport{
		Timestamp:       time.Now(),
		TotalRevenue:    totalRevenue,
		TotalOrders:     totalOrders,
		TopSellingBooks: topSellingBooks,
	}, nil
}


// generate the sales report every 24 hours

func StartSalesReportJob(
	ctx context.Context,
	reportService *SalesReportService,
) {
	ticker := time.NewTicker(24 * time.Hour)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Println("Starting daily sales report generation")

				to := time.Now()
				from := to.Add(-24 * time.Hour)

				report, err := reportService.GenerateSalesReport(ctx, from, to)
if err != nil {
	log.Printf("Failed to generate sales report: %v", err)
	continue
}

report.Timestamp = time.Now()

if err := SaveReportToJson(report); err != nil {
	log.Printf("Failed to save sales report: %v", err)
	continue
}

log.Printf(
	"Sales report SAVED: orders=%d revenue=%.2f",
	report.TotalOrders,
	report.TotalRevenue,
)


			case <-ctx.Done():
				log.Println("Stopping sales report job")
				return
			}
		}
	}()
}
