package main

import (
	"fmt"
	"net/http"

	"country-info-service/handlers"
)

func main() {
	// Register handlers
	http.HandleFunc("/countryinfo/v1/info/", handlers.CountryInfoHandler)
	http.HandleFunc("/countryinfo/v1/population/", handlers.PopulationHandler)
	http.HandleFunc("/countryinfo/v1/status/", handlers.StatusHandler)

	// Start server
	fmt.Println("Server is running on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
