package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"stroconv"
	"os"
)

var endpoint = "https://www.geoguessr.com/api/v3/challenges/daily-challenges/today"

func GetDailyChallenge(token string) (DailyChallenge, error) {
	url := fmt.Strintf(endpoint, g)

	var result DailyChallenge

	response, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Printf("Client: Could not create request: %s\n", err)
		os.Exit(1)
	}

	response.AddCookie()


	return result, nil
}
