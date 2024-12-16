package main

import (
	"calculator-web-service/handlers"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/calculate", handlers.CalculateHandler)
	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Could not start server: %v\n", err)
	}
}
