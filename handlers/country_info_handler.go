package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"fmt"
	"regexp"

	"country-info-service/utils"
)


// CountryInfoHandler handles requests for country information based on an ISO2 country code.
// It fetches country details and a list of major cities, with an optional limit on the number of cities.
//
// Endpoint: GET /countryinfo/v1/info/{code}?limit={number}
//
// Parameters:
//   - code: (string) The ISO2 country code (e.g., "no" for Norway).
//   - limit (optional): (int) The maximum number of cities to include in the response (default: 10).
//
// Example Requests:
//   - GET /countryinfo/v1/info/no
//   - GET /countryinfo/v1/info/us?limit=5
//
// Possible HTTP Status Codes:
//   - 200 OK: Request was successful.
//   - 400 Bad Request: Missing or invalid country code, or invalid query parameter.
//   - 500 Internal Server Error: Failed to fetch country information.
func CountryInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract country code from URL path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "missing country code. Example: /countryinfo/v1/info/no", http.StatusBadRequest)
		return
	}
	countryCode := strings.ToUpper(parts[4]) // Convert to uppercase (ISO2 codes are uppercase)
	// For debugging
	// fmt.Println("received country code:", countryCode)

	// Check if country code is ISO2 format
	if matched, _ := regexp.MatchString("^[A-Z]{2}$", countryCode); !matched {
		http.Error(w, "invalid country code. Use a valid ISO2 format (e.g., 'NO', 'US')", http.StatusBadRequest)
		return
	}

	// Extract the "limit" query parameter, defaulting to 10 if not provided
	limit := 10
	if queryLimit := r.URL.Query().Get("limit"); queryLimit != "" {
		parsedLimit, err := strconv.Atoi(queryLimit)
		if err != nil || parsedLimit <= 0 {
			http.Error(w, "invalid 'limit' parameter. Must be a positive integer.", http.StatusBadRequest)
			return
		}
		limit = parsedLimit
	}
	// for debugging
	// fmt.Println("Limit set to:", limit)

	// Fetch country information using the provided country code and limit
	info, err := utils.FetchCountryInfo(countryCode, limit)
	if err != nil {
		fmt.Println("Error fetching country info:", err)

		// Differentiate error types
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "country not found in the database.", http.StatusNotFound)
		} else if strings.Contains(err.Error(), "API request failed") {
			http.Error(w, "failed to retrieve data from the external API.", http.StatusBadGateway)
		} else {
			http.Error(w, "internal server error while fetching data.", http.StatusInternalServerError)
		}
		return
	}

	// Return the fetched country information as a JSON response
	// for debugging
	// fmt.Println("Returning country info:", info)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}
