package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Serve a simple HTML form at the root
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
            <form action="/weather" method="get">
                Address: <input type="text" name="address" />
                <input type="submit" value="Get Weather" />
            </form>
        `)
	})

	// Weather endpoint
	r.GET("/weather", func(c *gin.Context) {
		address := c.Query("address")

		if address == "" {
			c.String(http.StatusBadRequest, "Address parameter is missing")
			return
		}

		// Use OpenStreetMap Nominatim for geocoding
		baseURL := "https://nominatim.openstreetmap.org/search"
		params := url.Values{}
		params.Set("q", address)
		params.Set("format", "json")
		params.Set("limit", "1")

		fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
		req, _ := http.NewRequest("GET", fullURL, nil)
		req.Header.Set("User-Agent", "GoWeatherApp (your@email.com)")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to connect to geocoding API")
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read response")
			return
		}

		type NominatimResult struct {
			Lat         string `json:"lat"`
			Lon         string `json:"lon"`
			DisplayName string `json:"display_name"`
		}

		var results []NominatimResult
		if err := json.Unmarshal(body, &results); err != nil {
			c.String(http.StatusInternalServerError, "Failed to parse JSON")
			return
		}

		if len(results) == 0 {
			c.String(http.StatusOK, "No coordinates found for '%s'", address)
			return
		}

		lat := results[0].Lat
		lon := results[0].Lon

		pointsURL := fmt.Sprintf("https://api.weather.gov/points/%s,%s", lat, lon)
		pointsReq, _ := http.NewRequest("GET", pointsURL, nil)
		pointsReq.Header.Set("User-Agent", "GoWeatherApp (your@email.com)")
		pointsResp, err := http.DefaultClient.Do(pointsReq)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to fetch weather metadata")
			return
		}
		defer pointsResp.Body.Close()

		pointsBody, err := io.ReadAll(pointsResp.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read weather metadata")
			return
		}

		type PointsResponse struct {
			Properties struct {
				Forecast string `json:"forecast"`
			} `json:"properties"`
		}
		var points PointsResponse
		if err := json.Unmarshal(pointsBody, &points); err != nil {
			c.String(http.StatusInternalServerError, "Failed to parse weather metadata JSON")
			return
		}
		forecastURL := points.Properties.Forecast
		if forecastURL == "" {
			c.String(http.StatusOK, "Coordinates for '%s': Latitude: %s, Longitude: %s\nNo forecast URL found for this location.", address, lat, lon)
			return
		}

		forecastReq, _ := http.NewRequest("GET", forecastURL, nil)
		forecastReq.Header.Set("User-Agent", "GoWeatherApp (your@email.com)")
		forecastResp, err := http.DefaultClient.Do(forecastReq)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to fetch weather forecast")
			return
		}
		defer forecastResp.Body.Close()

		forecastBody, err := io.ReadAll(forecastResp.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to read weather forecast")
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
			c.String(http.StatusInternalServerError, "Failed to parse weather forecast JSON")
			return
		}

		if len(forecast.Properties.Periods) > 0 {
			p := forecast.Properties.Periods[0]
			c.String(http.StatusOK,
				"Coordinates for '%s': Latitude: %s, Longitude: %s\nWeather Forecast for %s: %d°F, %s",
				address, lat, lon, p.Name, p.Temperature, p.ShortForecast)
		} else {
			c.String(http.StatusOK,
				"Coordinates for '%s': Latitude: %s, Longitude: %s\nNo weather forecast data available.",
				address, lat, lon)
		}
	})

	r.Run(":8080")
}
