package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Miconen/geoguessr-daily-challenge-recap/models"
)

func GetEndpoint(lat float64, lng float64) string {
	return fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=%f&lon=%f&accept-language=en", lat, lng)
}

func NominatimRequest(ep string) (models.GeoData, error) {
	var result models.GeoData

	request, err := http.NewRequest("GET", ep, nil)
	if err != nil {
		fmt.Printf("Client: Could not create request: %s\n", err)
		os.Exit(1)
	}

	request.Header.Set("User-Agent", "GeoGuessrDailyRecapBot/1.0")

	// Actually execute the request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return result, fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()

	// Check status code
	if response.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	// Decode the JSON response
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
