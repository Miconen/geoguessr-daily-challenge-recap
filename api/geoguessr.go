package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Miconen/geoguessr-daily-challenge-recap/models"
)

var endpoint = "https://www.geoguessr.com/api/v3/challenges/daily-challenges/today"

func GetDailyChallenge(token string) (models.Challenge, error) {
	var result models.Challenge

	request, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Printf("Client: Could not create request: %s\n", err)
		os.Exit(1)
	}

	request.Header.Set("Cookie", fmt.Sprintf("_ncfa=%s", token))

	// Actually execute the request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return result, fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()

	// fmt.Printf("Status Code: %d\n", response.StatusCode)
	// fmt.Printf("Response Body: %s\n", string(body))

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
