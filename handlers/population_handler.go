package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"regexp"
	"time"

	"country-info-service/utils"
)

// PopulationHandler handles requests for country population data based on an ISO2 country code and optional year range.
//
// Endpoint: GET /countryinfo/v1/population/{countryCode}?limit={startYear-endYear}
//
// Parameters:
//   - countryCode: (string) The ISO2 country code (e.g., "NO" for Norway).
//   - limit (optional): (string) A year range in the format "startYear-endYear" (e.g., "2000-2020"). Has to be valid 4 digit year counts.
//
// Example Requests:
//   - GET /countryinfo/v1/population/NO
//   - GET /countryinfo/v1/population/US?limit=2000-2010
//
// Response:
//   A JSON object containing population data with mean value and an array of year-value pairs.
//
// Possible HTTP Status Codes:
//   - 200 OK: Request was successful.
//   - 400 Bad Request: Missing or invalid country code, or invalid query parameters.
//   - 404 Not Found: No population data available for the specified country or year range.
//   - 502 Bad Gateway: External API failure.
func PopulationHandler(w http.ResponseWriter, r *http.Request) {
	// Gets country code and validate it
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/countryinfo/v1/population/"), "/")
	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "Missing country code. Example: /countryinfo/v1/population/NO", http.StatusBadRequest)
		return
	}
	countryCode := strings.ToUpper(parts[0]) // Convert to uppercase for consistency and easier comparison

	// Check if the country code is in ISO2 format
	if matched, _ := regexp.MatchString("^[A-Z]{2}$", countryCode); !matched {
		http.Error(w, "Invalid country code. Use a valid ISO2 format (e.g., 'NO', 'US')", http.StatusBadRequest)
		return
	}

	// Parse optional limit query param (startYear-endYear)
	startYear, endYear := 0, 0
	limitParam := r.URL.Query().Get("limit")
	if limitParam != "" {
		years := strings.Split(limitParam, "-")
		if len(years) == 2 {
			var err1, err2 error
			startYear, err1 = strconv.Atoi(years[0])
			endYear, err2 = strconv.Atoi(years[1])

			// Validate year range
			if err1 != nil || err2 != nil {
				http.Error(w, "Invalid 'limit' format. Use 'startYear-endYear' with numeric values (e.g., '2000-2020').", http.StatusBadRequest)
				return
			}
			currentYear := time.Now().Year() // Gets current year

			if startYear < 1900 || endYear > currentYear {
				http.Error(w, fmt.Sprintf("Year range out of bounds. Use years between 1900 and %d.", currentYear), http.StatusBadRequest)
				return
			}
		} else {
			http.Error(w, "Invalid 'limit' format. Use 'startYear-endYear'.", http.StatusBadRequest)
			return
		}
	}

	// Debugging output
	// fmt.Printf("Request received for country: %s | Start year: %d | End year: %d\n", countryCode, startYear, endYear)

	// Fetch population data
	data, err := utils.FetchPopulationData(countryCode, startYear, endYear)
	if err != nil {
		fmt.Println("Error fetching population data:", err)

		// Differentiate error types
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Population data not found for the given country or year range.", http.StatusNotFound)
		} else if strings.Contains(err.Error(), "API request failed") {
			http.Error(w, "Failed to retrieve data from the external API.", http.StatusBadGateway)
		} else {
			http.Error(w, "Internal server error while fetching population data.", http.StatusInternalServerError)
		}
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
