package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GeoGuessrRequest[T any](ncfa string, ep string) (T, error) {
	var result T

	request, err := http.NewRequest("GET", ep, nil)
	if err != nil {
		fmt.Printf("Client: Could not create request: %s\n", err)
		os.Exit(1)
	}

	request.Header.Set("Cookie", fmt.Sprintf("_ncfa=%s", ncfa))
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// Actually execute the request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return result, fmt.Errorf("request failed: %w", err)
	}
	defer response.Body.Close()

	// Check status code
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return result, fmt.Errorf("unexpected status code: %d from URL %s (response: %s)", response.StatusCode, ep, string(body))
	}

	// Decode the JSON response
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
