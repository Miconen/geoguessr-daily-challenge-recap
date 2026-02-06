package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Miconen/geoguessr-daily-challenge-recap/api"
	"github.com/Miconen/geoguessr-daily-challenge-recap/discord"
)

var (
	DiscordToken   string
	GeoguessrToken string
	Users          string
)

func init() {
	flag.StringVar(&DiscordToken, "discord", "", "Bot Token")
	flag.StringVar(&GeoguessrToken, "geoguessr", "", "Geoguessr Token")
	flag.StringVar(&Users, "users", "", "Discord Users")
	flag.Parse()

	if DiscordToken != "" {
		os.Setenv("DISCORD_TOKEN", DiscordToken)
	}

	if GeoguessrToken != "" {
		os.Setenv("NCFA_TOKEN", GeoguessrToken)
	}

	if Users != "" {
		os.Setenv("DISCORD_USERS", Users)
	}
}

func main() {
	ncfa := os.Getenv("NCFA_TOKEN")
	if ncfa == "" {
		fmt.Fprintf(os.Stderr, "Error getting environment variable: \"NCFA_TOKEN\"\n")
		os.Exit(1)
	}

	dt := os.Getenv("DISCORD_TOKEN")
	if ncfa == "" {
		fmt.Fprintf(os.Stderr, "Error getting environment variable: \"DISCORD_TOKEN\"\n")
		os.Exit(1)
	}

	users := strings.Split(os.Getenv("DISCORD_USERS"), ",")

	challenge, err := api.GetDailyChallenge(ncfa)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching name changes: %v\n", err)
		os.Exit(1)
	}

	err = discord.SendDM(dt, users, challenge)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending Discord messages: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
