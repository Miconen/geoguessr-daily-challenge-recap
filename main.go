package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Miconen/geoguessr-daily-challenge-recap/api"
)

func main() {
	tokenFlag := flag.String("token", "", "token")

	flag.Parse()

	if *tokenFlag != "" {
		os.Setenv("TOKEN", *tokenFlag)
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		fmt.Fprintf(os.Stderr, "Error getting environment variable: \"TOKEN\"\n")
		os.Exit(1)
	}

	challenge, err := api.GetDailyChallenge(token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching name changes: %v\n", err)
		os.Exit(1)
	}

	for _, v := range challenge.Club {
		fmt.Println(v.Nick)
	}

	os.Exit(0)
}
