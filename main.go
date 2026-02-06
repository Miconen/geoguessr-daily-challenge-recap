package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Miconen/geoguessr-daily-challenge-recap/api"
	"github.com/Miconen/geoguessr-daily-challenge-recap/discord"
	"github.com/Miconen/geoguessr-daily-challenge-recap/models"
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

	challenge, err := api.GeoGuessrRequest[models.Challenge](ncfa, api.EndpointDaily)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching name changes: %v\n", err)
		os.Exit(1)
	}

	competition, err := api.GeoGuessrRequest[models.Competition](ncfa, api.GetScoresEndpoint(challenge.Token))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching daily challenge: %v\n", err)
		os.Exit(1)
	}

	if len(competition.Items) == 0 {
		fmt.Fprintf(os.Stderr, "Error fetching games by challenge id: %v\n", err)
		os.Exit(1)
	}

	GeoData := models.GameGeoData{
		ActualLocations: []models.RoundGeoData{},
		PlayerGuesses:   []models.PlayerGeoData{},
	}

	// 1. Process Actual Locations (Only once from the first player)
	for i, round := range competition.Items[0].Game.Rounds {
		// These are the "Correct" answers
		ep := api.GetEndpoint(round.Lat, round.Lng)

		location, err := api.NominatimRequest(ep)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching actual location R%d: %v\n", i+1, err)
			continue
		}

		// Store in your GeoData.ActualLocations slice
		fmt.Printf("Round %d is in: %s\n", i+1, location.Address.Country)
		data := models.RoundGeoData{
			RoundNumber: strconv.Itoa(i + 1),
			Location:    location,
		}
		GeoData.ActualLocations = append(GeoData.ActualLocations, data)

		time.Sleep(time.Second)
	}

	// 2. Process Player Guesses
	for _, item := range competition.Items {
		fmt.Printf("Processing guesses for: %s\n", item.Game.Player.Nick)
		player := models.PlayerGeoData{
			PlayerID: item.Game.Player.ID,
			Rounds:   make([]models.GuessGeoData, 0),
		}

		for i, guess := range item.Game.Player.Guesses {
			// These are the coordinates where the player clicked
			ep := api.GetEndpoint(guess.Lat, guess.Lng)

			location, err := api.NominatimRequest(ep)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching guess %d for %s: %v\n", i+1, item.Game.Player.Nick, err)
				continue
			}

			// Store in your GeoData.PlayerGuesses
			fmt.Printf("  - Guess %d: %s\n", i+1, location.Address.Country)
			guess := models.GuessGeoData{
				RoundNumber: i + 1,
				Guess:       location,
			}
			player.Rounds = append(player.Rounds, guess)

			time.Sleep(time.Second)
		}
		GeoData.PlayerGuesses = append(GeoData.PlayerGuesses, player)
	}

	embed := discord.GenerateGeoGuessrDailyChallengeEmbed(competition.Items, challenge, &GeoData)

	err = discord.SendDM(dt, users, embed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error sending Discord messages: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
