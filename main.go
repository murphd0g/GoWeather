package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

func main() {
	http.HandleFunc("/", homePage)
	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	benchmark := r.URL.Query().Get("benchmark")
	format := r.URL.Query().Get("format")

	if address == "" {
		http.Error(w, "Address parameter is missing", http.StatusBadRequest)
		return
	}
	if benchmark == "" {
		benchmark = "2020" // default value
	}
	if format == "" {
		format = "json" // default value
	}

	baseURL := "https://geocoding.geo.census.gov/geocoder/locations/onelineaddress"
	params := url.Values{}
	params.Set("address", address)
	params.Set("benchmark", benchmark)
	params.Set("format", format)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		http.Error(w, "Failed to connect to geocoding API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}

	type GeocodeResult struct {
		Result struct {
			AddressMatches []struct {
				Coordinates struct {
					X float64 `json:"x"`
					Y float64 `json:"y"`
				} `json:"coordinates"`
			} `json:"addressMatches"`
		} `json:"result"`
	}

	var result GeocodeResult
	if err := json.Unmarshal(body, &result); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	if len(result.Result.AddressMatches) > 0 {
		coords := result.Result.AddressMatches[0].Coordinates
		fmt.Fprintf(w, "Coordinates for '%s': Latitude: %f, Longitude: %f", address, coords.Y, coords.X)
	} else {
		fmt.Fprintf(w, "No coordinates found for '%s'", address)
	}
}
