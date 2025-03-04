package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var cost float64

// Simulates a call to an external API that returns the price of a product
func fetchProductPrice(productID string) (float64, error) {
	log.Printf("[COST: $0.01] Calling external API for product: %s\n", productID)
	time.Sleep(2 * time.Second) // Simulates latency
	cost += 0.01
	return 99.99, nil // Fictitious price
}

// Handler to fetch the price of a product
func getProductPriceHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the product ID from the URL
	productID := r.URL.Query().Get("id")

	// Fetches the product price from the external API
	price, err := fetchProductPrice(productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a response map
	response := map[string]interface{}{
		"product_id": productID,
		"price":      price,
	}

	// Convert the response to JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getCost(w http.ResponseWriter, r *http.Request) {
	// Create a response map
	response := map[string]interface{}{
		"total_cost": fmt.Sprintf("%.2f", cost),
	}
	// Convert the response to JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func clearCosts(w http.ResponseWriter, r *http.Request) {
	cost = 0
}

func main() {
	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register the handler for the /products/:id/price route
	mux.HandleFunc("/products/{id}/price", getProductPriceHandler)
	mux.HandleFunc("/costs", getCost)
	mux.HandleFunc("/clear-costs", clearCosts)

	// Start the server
	log.Println("API running without singleflight on port :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
