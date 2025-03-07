package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// The response structure for population data.
type PopulationResponse struct {
	Mean   int `json:"mean"`
	Values []struct {
		Year  int `json:"year"`
		Value int `json:"value"`
	} `json:"values"`
}

// The response structure for population API data.
type ApiResponse struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
	Data  struct {
		Country          string `json:"country"`
		Code             string `json:"code"`
		Iso3             string `json:"iso3"`
		PopulationCounts []struct {
			Year  int `json:"year"`
			Value int `json:"value"`
		} `json:"populationCounts"`
	} `json:"data"`
}

// Response structure to hold the common name of a country.
type CountryNameResponse []struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
}

// FetchCountryName retrieves the common name of a country using its ISO2 code.
func FetchCountryName(iso2 string) (string, error) {
	url := fmt.Sprintf("http://129.241.150.113:8080/v3.1/alpha/%s", iso2)

	// Fetch country name from the API
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching country name: %v", err)
		return "", fmt.Errorf("failed to fetch country name: %w", err)
	}
	defer resp.Body.Close()

	// Check if the API returned a valid response
	if resp.StatusCode != http.StatusOK {
		log.Printf("Country API returned status %d", resp.StatusCode)
		return "", fmt.Errorf("country API returned status %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading country API response: %v", err)
		return "", fmt.Errorf("failed to read country API response: %w", err)
	}

	// Decode the response body
	var countryData CountryNameResponse
	if err := json.Unmarshal(body, &countryData); err != nil {
		log.Printf("Error decoding country API response: %v", err)
		return "", fmt.Errorf("failed to decode country API response: %w", err)
	}

	// Ensure response is not empty
	if len(countryData) == 0 || countryData[0].Name.Common == "" {
		log.Printf("Invalid country name received for ISO2 code: %s", iso2)
		return "", fmt.Errorf("invalid country name received for ISO2 code: %s", iso2)
	}

	return countryData[0].Name.Common, nil
}

// FetchPopulationData retrieves population data for a country within a given year range.
func FetchPopulationData(iso2 string, startYear, endYear int) (*PopulationResponse, error) {
	// Validate inputs
	if startYear > endYear && endYear != 0 {
		return nil, fmt.Errorf("invalid year range: startYear (%d) cannot be greater than endYear (%d)", startYear, endYear)
	}

	// Fetch country name
	countryName, err := FetchCountryName(iso2)
	if err != nil {
		return nil, fmt.Errorf("failed to get country name: %w", err)
	}

	// Prepare the JSON for the POST request
	apiRequest := struct {
		Country string `json:"country"`
	}{
		Country: countryName,
	}

	// Marshal the JSON data
	jsonData, err := json.Marshal(apiRequest)
	if err != nil {
		log.Printf("Error marshaling JSON request: %v", err)
		return nil, fmt.Errorf("failed to create JSON request: %w", err)
	}

	// Send POST request to the population API
	populationAPI := "http://129.241.150.113:3500/api/v0.1/countries/population"
	resp, err := http.Post(populationAPI, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error sending request to Population API: %v", err)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check API response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Population API returned status: %d", resp.StatusCode)
		return nil, fmt.Errorf("population API returned status %d", resp.StatusCode)
	}

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading Population API response: %v", err)
		return nil, fmt.Errorf("failed to read Population API response: %w", err)
	}

	// Decode API response
	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Error decoding Population API response: %v", err)
		return nil, fmt.Errorf("failed to decode Population API response: %w", err)
	}

	// Check for errors in the API response
	if apiResponse.Error {
		log.Printf("Population API error: %s", apiResponse.Msg)
		return nil, fmt.Errorf("population API returned an error: %s", apiResponse.Msg)
	}

	// Extract population data
	populationData := apiResponse.Data.PopulationCounts
	if populationData == nil {
		log.Printf("No population data found for country: %s", countryName)
		return nil, fmt.Errorf("no population data found for country: %s", countryName)
	}

	// FIlter population data based on year range
	var filteredCounts []struct {
		Year  int `json:"year"`
		Value int `json:"value"`
	}
	total, count := 0, 0
	for _, entry := range populationData {
		if (startYear == 0 || entry.Year >= startYear) && (endYear == 0 || entry.Year <= endYear) {
			filteredCounts = append(filteredCounts, entry)
			total += entry.Value
			count++
		}
	}

	// Calculate mean population
	mean := 0
	if count > 0 {
		mean = total / count
	}

	// Return the response
	return &PopulationResponse{
		Mean:   mean,
		Values: filteredCounts,
	}, nil
}
