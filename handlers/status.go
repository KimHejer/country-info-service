package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"log"
)

// The status handler provides real-time service diagnostics from the API endpoints used in the country-info-service.
// It checks the health status of the CountriesNow API and the RestCountries API.
// The uptime of the service is also calculated and returned in the response.
//
// Endpoint: GET /countryinfo/v1/status
//
// Example Request:
//   GET /countryinfo/v1/status
//
// Response:
//   A JSON object containing the health status of the APIs, the current version and the service uptime in seconds.
//
// Possible HTTP Status Codes:
//   - 200 OK: Request was successful.
//   - 502 Bad Gateway: External API failure.
//
// Example Response:
//   {
//     "countriesnowapi": "200",
//     "restcountriesapi": "200",
//     "version": "v1",
//     "uptime": 128
//   }


// Start time for uptime tracking
var startTime = time.Now()

// APIStatus represents the health status of an API
type APIStatus struct {
	CountriesNowAPI  string `json:"countriesnowapi"`  // Status of the CountriesNow API
	RestCountriesAPI string `json:"restcountriesapi"` // Status of the RestCountries API
	Version          string `json:"version"`          // API version
	Uptime           int    `json:"uptime"`           // Service uptime in seconds
}

// checkAPIHealth makes a request to an API with a timeout and returns its status
func checkAPIHealth(url string) string {
	client := &http.Client{Timeout: 3 * time.Second} // Set a 3-second timeout

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error checking API %s: %v\n", url, err)
		return "FAILED"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API %s returned status %d\n", url, resp.StatusCode)
		return "ERROR " + http.StatusText(resp.StatusCode)
	}

	return "200"
}

// StatusHandler provides real-time service diagnostics
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	uptime := int(time.Since(startTime).Seconds())

	// Check both APIs
	countriesNowStatus := checkAPIHealth("http://129.241.150.113:3500/api/v0.1/countries")
	restCountriesStatus := checkAPIHealth("http://129.241.150.113:8080/v3.1/all")

	// Construct JSON response
	status := APIStatus{
		CountriesNowAPI:  countriesNowStatus,
		RestCountriesAPI: restCountriesStatus,
		Version:          "v1",
		Uptime:           uptime,
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
