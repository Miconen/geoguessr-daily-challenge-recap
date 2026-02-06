package main

import (
	"context"
	"flag"
	"fmt"
	"os"
)

func main() {
	friendsFlag := flag.String("friends", "", "friends")
	tokenFlag := flag.String("token", "", "token")

	flag.Parse()

	if *friendsFlag != "" {
		os.Setenv("FRIENDS", *friendsFlag)
	}
	if *tokenFlag != "" {
		os.Setenv("TOKEN", *tokenFlag)
	}

	token := os.Getend("TOKEN")
	if token == "" {
		fmt.Fprintf(os.Stderr, "Error getting environment variable: \"TOKEN\"\n")
		os.Exit(1)
	}

	friends := os.Getenv("FRIENDS")

	users, err := utils.(group)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching name changes: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
