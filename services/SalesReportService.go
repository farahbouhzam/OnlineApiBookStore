package services

import (
	"context"
	"encoding/json"
	"fmt"
	"online_bookStore/Interfaces"
	"online_bookStore/models"
	"path/filepath"
     "os"
	"log"
	"sort"
	"time"
)


type SalesReportService struct {
	orderStore interfaces.OrderStore
}

func SaveReportToJson(report models.SalesReport) error {
	
	// filename based on date

	filename := fmt.Sprintf(
		"sales_report_%s.json", time.Now().Format("2006-01-02"),
	)

	path := filepath.Join("reports", filename)

	// create file
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	encoder.SetIndent(""," ")
	return encoder.Encode(report)
}

// Constructeur
func NewSalesReportService(orderStore interfaces.OrderStore) *SalesReportService{
	return &SalesReportService{
		orderStore : orderStore,
	}
}

func (s *SalesReportService) GenerateSalesReport(ctx context.Context,from time.Time,to time.Time)(models.SalesReport, error){

	orders, err := s.orderStore.GetOrderByDateRange(ctx,from,to)
    
	TotalRevenue := 0.0
	totalOrders := len(orders)

	if err != nil {
		return models.SalesReport{}, err
	}


	bookCounter := make(map[int]int)
	

	
	for _,order := range orders {
		TotalRevenue += order.TotalPrice

        for _, item := range order.Items {
			bookCounter[item.ID]+= item.Quantity
		}

	}

	var sales []models.BookSales

	for key, value := range bookCounter {
		sales = append(sales , models.BookSales{
			BookID: key,
			Quantity: value,

		})
	}

	sort.Slice(sales, func(i,j int)bool{
          return sales[i].Quantity> sales[j].Quantity
	})

	var TopSellingBooks []models.BookSales
    limit := 3
	if len(sales) < limit {
		limit = len(sales)
	}
	for i:=0; i<limit; i++ {
		TopSellingBooks = append(TopSellingBooks, sales[i])
	}


	return models.SalesReport{
          TotalRevenue: TotalRevenue,
          TotalOrders: totalOrders,
		  TopSellingBooks: TopSellingBooks,
	}, nil



}


// generate the sales report every 24 hours

func StartSalesReportJob(ctx context.Context, reportService *SalesReportService){

	ticker := time.NewTicker(24*time.Hour)

	go func(){
		// Stops the ticker when the goroutine exits
		defer ticker.Stop()
  
		for {
			select {
			case <- ticker.C :
				log.Println("Starting daily sales report generation")

				to := time.Now()
				from := to.Add(-24*time.Hour)


				report, err := reportService.GenerateSalesReport(ctx,from,to)
				if err != nil {
					log.Printf("Failed to generate sales report: %v", err)
					continue
				}

				log.Printf(
					"Sales report generated: orders=%d revenue=%.2f",
					report.TotalOrders,
					report.TotalRevenue,
				)
			case <-ctx.Done():
				log.Println("Stopping sales report job")
				return
			}
			} }() }


		







