package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
)

// CountryInfoResponse represents the structured response for country information.
type CountryInfoResponse struct {
	Name       string            `json:"name"`       
	Continent  string            `json:"continent"`  
	Population int               `json:"population"` 
	Languages  map[string]string `json:"languages"`  
	Borders    []string          `json:"borders"`   
	Flag       string            `json:"flag"`       
	Capital    string            `json:"capital"`   
	Cities     []string          `json:"cities"`   
}

// FetchCountryInfo queries the REST Countries API and the Cities API to get country details.
func FetchCountryInfo(countryCode string, limit int) (*CountryInfoResponse, error) {
	url := fmt.Sprintf("http://129.241.150.113:8080/v3.1/alpha/%s", countryCode)

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching country data: %v", err)
		return nil, fmt.Errorf("failed to fetch country data: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("REST Countries API returned status: %d", resp.StatusCode)
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Decode JSON response
	var data []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Error decoding country API response: %v", err)
		return nil, fmt.Errorf("failed to decode country API response: %w", err)
	}

	// Ensure response contains data
	if len(data) == 0 {
		log.Printf("No data found for country code: %s", countryCode)
		return nil, fmt.Errorf("no data found for country code: %s", countryCode)
	}

	country := data[0]

	// Extract country name
	name, ok := extractString(country, "name", "common")
	if !ok {
		log.Printf("Country name not found in API response")
		return nil, fmt.Errorf("country name not found")
	}

	// Extract region
	region, ok := country["region"].(string)
	if !ok {
		region = "Unknown"
	}

	// Extract borders
	borders := extractStringArray(country, "borders")

	// Extract languages
	languages := make(map[string]string)
	if langs, ok := country["languages"].(map[string]interface{}); ok {
		for abbr, lang := range langs {
			if langName, valid := lang.(string); valid {
				languages[abbr] = langName
			}
		}
	}

	// Extract flag URL
	flag, _ := extractString(country, "flags", "svg")

	// Extract capital city
	capital := "N/A"
	if capList, ok := country["capital"].([]interface{}); ok && len(capList) > 0 {
		if capStr, valid := capList[0].(string); valid {
			capital = capStr
		}
	}

	// Extract population
	population := 0
	if pop, ok := country["population"].(float64); ok {
		population = int(pop)
	}

	// Fetch cities
	cities, err := FetchCities(name, limit)
	if err != nil {
		log.Printf("Error fetching cities: %v", err)
		cities = []string{"City data not available"}
	}

	// Construct the response
	response := CountryInfoResponse{
		Name:       name,
		Continent:  region,
		Population: population,
		Languages:  languages,
		Borders:    borders,
		Flag:       flag,
		Capital:    capital,
		Cities:     cities,
	}

	return &response, nil
}

// FetchCities queries the Cities API to get a list of cities for a given country.
func FetchCities(countryName string, limit int) ([]string, error) {
	url := "http://129.241.150.113:3500/api/v0.1/countries/cities"
	payload := []byte(fmt.Sprintf(`{"country": "%s"}`, countryName))

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Error creating cities API request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request to Cities API: %v", err)
		return nil, fmt.Errorf("failed to fetch cities: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP response status
	if resp.StatusCode != http.StatusOK {
		log.Printf("Cities API returned status: %d", resp.StatusCode)
		return nil, fmt.Errorf("cities API returned status %d", resp.StatusCode)
	}

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading cities API response: %v", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse API response
	var apiResponse struct {
		Error  bool     `json:"error"`
		Cities []string `json:"data"`
	}
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		log.Printf("Error decoding cities API response: %v", err)
		return nil, fmt.Errorf("failed to parse cities API response: %w", err)
	}

	// Handle API errors
	if apiResponse.Error {
		log.Printf("Cities API reported an error for country: %s", countryName)
		return nil, fmt.Errorf("error fetching cities for country: %s", countryName)
	}

	// Apply limit to cities
	if len(apiResponse.Cities) > limit {
		return apiResponse.Cities[:limit], nil
	}

	return apiResponse.Cities, nil
}

// Extracts a nested string value from a map.
func extractString(data map[string]interface{}, keys ...string) (string, bool) {
	for _, key := range keys {
		if nestedMap, ok := data[key].(map[string]interface{}); ok {
			for _, nestedKey := range keys {
				if value, exists := nestedMap[nestedKey]; exists {
					if strValue, valid := value.(string); valid {
						return strValue, true
					}
				}
			}
		} else if value, exists := data[key]; exists {
			if strValue, valid := value.(string); valid {
				return strValue, true
			}
		}
	}
	return "", false
}

// Extracts an array of strings from a map.
func extractStringArray(data map[string]interface{}, key string) []string {
	result := []string{}
	if arr, ok := data[key].([]interface{}); ok {
		for _, item := range arr {
			if strItem, valid := item.(string); valid {
				result = append(result, strItem)
			}
		}
	}
	return result
}
