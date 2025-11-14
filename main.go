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
		lat := fmt.Sprintf("%.4f", coords.Y)
		lon := fmt.Sprintf("%.4f", coords.X)

		// Step 1: Get metadata for the location
		pointsURL := fmt.Sprintf("https://api.weather.gov/points/%s,%s", lat, lon)
		req, _ := http.NewRequest("GET", pointsURL, nil)
		req.Header.Set("User-Agent", "GoWeatherApp (your@email.com)")
		pointsResp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "Failed to fetch weather metadata", http.StatusInternalServerError)
			return
		}
		defer pointsResp.Body.Close()

		pointsBody, err := io.ReadAll(pointsResp.Body)
		if err != nil {
			http.Error(w, "Failed to read weather metadata", http.StatusInternalServerError)
			return
		}

		// Step 2: Extract the forecast URL from the metadata
		type PointsResponse struct {
			Properties struct {
				Forecast string `json:"forecast"`
			} `json:"properties"`
		}
		var points PointsResponse
		if err := json.Unmarshal(pointsBody, &points); err != nil {
			http.Error(w, "Failed to parse weather metadata JSON", http.StatusInternalServerError)
			return
		}
		forecastURL := points.Properties.Forecast
		if forecastURL == "" {
			fmt.Fprintf(w, "Coordinates for '%s': Latitude: %.4f, Longitude: %.4f\n", address, coords.Y, coords.X)
			fmt.Fprintf(w, "No forecast URL found for this location.")
			return
		}

		// Step 3: Query the forecast URL
		req2, _ := http.NewRequest("GET", forecastURL, nil)
		req2.Header.Set("User-Agent", "GoWeatherApp (your@email.com)")
		forecastResp, err := http.DefaultClient.Do(req2)
		if err != nil {
			http.Error(w, "Failed to fetch weather forecast", http.StatusInternalServerError)
			return
		}
		defer forecastResp.Body.Close()

		forecastBody, err := io.ReadAll(forecastResp.Body)
		if err != nil {
			http.Error(w, "Failed to read weather forecast", http.StatusInternalServerError)
			return
		}

		type ForecastResponse struct {
			Properties struct {
				Periods []struct {
					Name          string `json:"name"`
					Temperature   int    `json:"temperature"`
					ShortForecast string `json:"shortForecast"`
				} `json:"periods"`
			} `json:"properties"`
		}
		var forecast ForecastResponse
		if err := json.Unmarshal(forecastBody, &forecast); err != nil {
			http.Error(w, "Failed to parse weather forecast JSON", http.StatusInternalServerError)
			return
		}

		if len(forecast.Properties.Periods) > 0 {
			p := forecast.Properties.Periods[0]
			fmt.Fprintf(w, "Coordinates for '%s': Latitude: %.4f, Longitude: %.4f\n", address, coords.Y, coords.X)
			fmt.Fprintf(w, "Weather Forecast for %s: %d°F, %s", p.Name, p.Temperature, p.ShortForecast)
		} else {
			fmt.Fprintf(w, "Coordinates for '%s': Latitude: %.4f, Longitude: %.4f\n", address, coords.Y, coords.X)
			fmt.Fprintf(w, "No weather forecast data available.")
		}
	} else {
		fmt.Fprintf(w, "No coordinates found for '%s'", address)
	}
}
