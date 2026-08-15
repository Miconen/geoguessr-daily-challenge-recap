package discord

import (
	"testing"
	"time"

	"github.com/Miconen/geoguessr-daily-challenge-recap/models"
)

func TestGenerateGeoGuessrDailyChallengeEmbed(t *testing.T) {
	challenge := models.Challenge{
		Date: time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC),
		Club: []models.ClubPlayer{
			{ID: "p1", Nick: "Mico", CurrentStreak: 44},
			{ID: "p2", Nick: "Cerkie", CurrentStreak: 54},
			{ID: "p3", Nick: "Anna", CurrentStreak: 12},
		},
	}

	items := []models.Items{
		{
			Game: models.Game{
				Player: models.Player{
					ID:          "p1",
					Nick:        "Mico",
					CountryCode: "fi",
					TotalScore:  models.TotalScore{Amount: "23154"},
					TotalDistanceInMeters: 608400,
					TotalTime:   268,
					Guesses: []models.Guesses{
						{RoundScoreInPoints: 4997, DistanceInMeters: 12000, Time: 180},
						{RoundScoreInPoints: 5000, DistanceInMeters: 14, Time: 36},
						{RoundScoreInPoints: 4340, DistanceInMeters: 142000, Time: 18},
						{RoundScoreInPoints: 3888, DistanceInMeters: 240000, Time: 23},
						{RoundScoreInPoints: 4929, DistanceInMeters: 35000, Time: 12},
					},
				},
			},
		},
		{
			Game: models.Game{
				Player: models.Player{
					ID:          "p2",
					Nick:        "Cerkie",
					CountryCode: "dk",
					TotalScore:  models.TotalScore{Amount: "21488"},
					TotalDistanceInMeters: 1272000,
					TotalTime:   480,
					Guesses: []models.Guesses{
						{RoundScoreInPoints: 4986, DistanceInMeters: 18000, Time: 180},
						{RoundScoreInPoints: 5000, DistanceInMeters: 42, Time: 154},
						{RoundScoreInPoints: 3239, DistanceInMeters: 512000, Time: 58},
						{RoundScoreInPoints: 3376, DistanceInMeters: 380000, Time: 26},
						{RoundScoreInPoints: 4887, DistanceInMeters: 48000, Time: 61},
					},
				},
			},
		},
		{
			Game: models.Game{
				Player: models.Player{
					ID:          "p3",
					Nick:        "Anna",
					CountryCode: "se",
					TotalScore:  models.TotalScore{Amount: "19850"},
					TotalDistanceInMeters: 2410000,
					TotalTime:   375,
					Guesses: []models.Guesses{
						{RoundScoreInPoints: 4810, DistanceInMeters: 54000, Time: 100},
						{RoundScoreInPoints: 4992, DistanceInMeters: 110, Time: 75},
						{RoundScoreInPoints: 3850, DistanceInMeters: 290000, Time: 62},
						{RoundScoreInPoints: 4120, DistanceInMeters: 180000, Time: 45},
						{RoundScoreInPoints: 4670, DistanceInMeters: 95000, Time: 55},
					},
				},
			},
		},
	}

	geodata := &models.GameGeoData{
		ActualLocations: []models.RoundGeoData{
			{Location: models.GeoData{Address: models.Address{Country: "Australia", State: "Queensland", CountryCode: "au"}}},
			{Location: models.GeoData{Address: models.Address{Country: "Finland", State: "South Savo", CountryCode: "fi"}}},
			{Location: models.GeoData{Address: models.Address{Country: "South Africa", State: "Western Cape", CountryCode: "za"}}},
			{Location: models.GeoData{Address: models.Address{Country: "Ukraine", State: "Kyiv Oblast", CountryCode: "ua"}}},
			{Location: models.GeoData{Address: models.Address{Country: "United Arab Emirates", State: "Sharjah Emirate", CountryCode: "ae"}}},
		},
		PlayerGuesses: []models.PlayerGeoData{
			{
				PlayerID: "p2",
				Rounds: []models.GuessGeoData{
					{RoundNumber: 1, Guess: models.GeoData{Address: models.Address{Country: "Australia", State: "Queensland", CountryCode: "au"}}},
					{RoundNumber: 2, Guess: models.GeoData{Address: models.Address{Country: "Finland", State: "South Savo", CountryCode: "fi"}}},
					{RoundNumber: 3, Guess: models.GeoData{Address: models.Address{Country: "Namibia", State: "Karas", CountryCode: "na"}}},
					{RoundNumber: 4, Guess: models.GeoData{Address: models.Address{Country: "Ukraine", State: "Kyiv Oblast", CountryCode: "ua"}}},
					{RoundNumber: 5, Guess: models.GeoData{Address: models.Address{Country: "United Arab Emirates", State: "Sharjah Emirate", CountryCode: "ae"}}},
				},
			},
		},
	}

	embed := GenerateGeoGuessrDailyChallengeEmbed(items, challenge, geodata)

	if embed == nil {
		t.Fatal("expected embed to be non-nil")
	}

	if len(embed.Fields) < 3 {
		t.Fatalf("expected at least 3 fields (Leaderboard, Rounds, Highlights), got %d", len(embed.Fields))
	}

	for _, field := range embed.Fields {
		t.Logf("\n--- %s ---\n%s\n", field.Name, field.Value)
	}
}
