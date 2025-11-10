package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"service/internal/shared/config"
	"time"
)

// GoogleMapsService handles Google Maps API integration
type GoogleMapsService struct {
	apiKey     string
	httpClient *http.Client
}

// NewGoogleMapsService creates a new Google Maps service
func NewGoogleMapsService() *GoogleMapsService {
	return &GoogleMapsService{
		apiKey: config.Config.GoogleMapsAPIKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// DistanceMatrixResponse represents the response from Google Maps Distance Matrix API
type DistanceMatrixResponse struct {
	Status            string   `json:"status"`
	OriginAddresses   []string `json:"origin_addresses"`
	DestinationAddresses []string `json:"destination_addresses"`
	Rows              []struct {
		Elements []struct {
			Status string `json:"status"`
			Distance struct {
				Value int    `json:"value"` // Distance in meters
				Text  string `json:"text"`
			} `json:"distance"`
			Duration struct {
				Value int    `json:"value"` // Duration in seconds
				Text  string `json:"text"`
			} `json:"duration"`
			DurationInTraffic struct {
				Value int    `json:"value"` // Duration in traffic in seconds
				Text  string `json:"text"`
			} `json:"duration_in_traffic"`
		} `json:"elements"`
	} `json:"rows"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// RouteInfo contains distance and duration information from Google Maps
type RouteInfo struct {
	Distance        float64 // Distance in kilometers
	Duration        int     // Duration in minutes
	DurationInTraffic int   // Duration in traffic in minutes (if available)
	Status          string  // Status of the API call
}

// GetDistanceAndDuration calculates distance and duration between two points using Google Maps Distance Matrix API
func (g *GoogleMapsService) GetDistanceAndDuration(originLat, originLon, destLat, destLon float64) (*RouteInfo, error) {
	// If API key is not configured, return error to fallback to Haversine
	if g.apiKey == "" {
		log.Println("[WARN] Google Maps API key not configured, will use Haversine fallback")
		return nil, errors.New("Google Maps API key not configured")
	}

	// Build the API URL
	baseURL := "https://maps.googleapis.com/maps/api/distancematrix/json"
	params := url.Values{}
	params.Add("origins", fmt.Sprintf("%f,%f", originLat, originLon))
	params.Add("destinations", fmt.Sprintf("%f,%f", destLat, destLon))
	params.Add("key", g.apiKey)
	params.Add("units", "metric") // Use metric units (kilometers)
	params.Add("mode", "driving") // Use driving mode
	params.Add("departure_time", "now") // Get real-time traffic data

	apiURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Make the HTTP request
	log.Printf("[INFO] Calling Google Maps Distance Matrix API: origin=(%f,%f), destination=(%f,%f)", originLat, originLon, destLat, destLon)
	resp, err := g.httpClient.Get(apiURL)
	if err != nil {
		log.Printf("[ERROR] Failed to call Google Maps API: %v", err)
		return nil, fmt.Errorf("failed to call Google Maps API: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the JSON response
	var apiResp DistanceMatrixResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if apiResp.Status != "OK" {
		log.Printf("[ERROR] Google Maps API error: %s - %s", apiResp.Status, apiResp.ErrorMessage)
		return nil, fmt.Errorf("Google Maps API error: %s - %s", apiResp.Status, apiResp.ErrorMessage)
	}

	// Extract distance and duration from response
	if len(apiResp.Rows) == 0 || len(apiResp.Rows[0].Elements) == 0 {
		return nil, errors.New("no route information found in response")
	}

	element := apiResp.Rows[0].Elements[0]
	if element.Status != "OK" {
		return nil, fmt.Errorf("route calculation failed: %s", element.Status)
	}

	// Convert distance from meters to kilometers
	distanceKm := float64(element.Distance.Value) / 1000.0

	// Convert duration from seconds to minutes
	durationMinutes := element.Duration.Value / 60

	// Get duration in traffic if available
	durationInTrafficMinutes := durationMinutes
	if element.DurationInTraffic.Value > 0 {
		durationInTrafficMinutes = element.DurationInTraffic.Value / 60
		log.Printf("[INFO] Google Maps API success: distance=%.2f km, duration=%d min, duration_in_traffic=%d min", distanceKm, durationMinutes, durationInTrafficMinutes)
	} else {
		log.Printf("[INFO] Google Maps API success: distance=%.2f km, duration=%d min (no traffic data)", distanceKm, durationMinutes)
	}

	return &RouteInfo{
		Distance:          distanceKm,
		Duration:          durationMinutes,
		DurationInTraffic: durationInTrafficMinutes,
		Status:            "OK",
	}, nil
}

// IsConfigured checks if Google Maps API is configured
func (g *GoogleMapsService) IsConfigured() bool {
	return g.apiKey != ""
}

