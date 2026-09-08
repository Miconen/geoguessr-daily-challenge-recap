package discord

import (
	"fmt"
	"testing"
	"time"

	"github.com/Miconen/geoguessr-daily-challenge-recap/models"
	"github.com/bwmarrin/discordgo"
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
					ID:                    "p1",
					Nick:                  "Mico",
					CountryCode:           "fi",
					TotalScore:            models.TotalScore{Amount: "23154"},
					TotalDistanceInMeters: 608400,
					TotalTime:             268,
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
					ID:                    "p2",
					Nick:                  "Cerkie",
					CountryCode:           "dk",
					TotalScore:            models.TotalScore{Amount: "21488"},
					TotalDistanceInMeters: 1272000,
					TotalTime:             480,
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
					ID:                    "p3",
					Nick:                  "Anna",
					CountryCode:           "se",
					TotalScore:            models.TotalScore{Amount: "19850"},
					TotalDistanceInMeters: 2410000,
					TotalTime:             375,
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

	if len(embed.Fields) != 7 {
		t.Fatalf("expected 7 fields (Leaderboard, 5 rounds, Highlights), got %d", len(embed.Fields))
	}

	assertEmbedLimits(t, embed)

	for _, field := range embed.Fields {
		t.Logf("\n--- %s ---\n%s\n", field.Name, field.Value)
	}
}

// TestEmbedLimitsWithManyPlayers reproduces the production 400 "Must be 1024 or fewer in length"
// by generating a full 26-player club and checking every Discord embed limit.
func TestEmbedLimitsWithManyPlayers(t *testing.T) {
	const players = 26
	challenge := models.Challenge{Date: time.Now()}
	var items []models.Items
	var guesses []models.PlayerGeoData
	for i := 0; i < players; i++ {
		id := fmt.Sprintf("p%02d", i)
		nick := fmt.Sprintf("PlayerWithALongName%02d", i)
		challenge.Club = append(challenge.Club, models.ClubPlayer{ID: id, Nick: nick, CurrentStreak: 100 + i})
		var gs []models.Guesses
		var geo []models.GuessGeoData
		for r := 0; r < 5; r++ {
			gs = append(gs, models.Guesses{RoundScoreInPoints: 4000 + i*10 + r, DistanceInMeters: 123456, Time: 125})
			geo = append(geo, models.GuessGeoData{RoundNumber: r + 1, Guess: models.GeoData{Address: models.Address{Country: "Brazil", CountryCode: "br"}}})
		}
		items = append(items, models.Items{Game: models.Game{Player: models.Player{
			ID: id, Nick: nick, CountryCode: "fi",
			TotalScore:            models.TotalScore{Amount: "20000"},
			TotalDistanceInMeters: 1234567, TotalTime: 3600, Guesses: gs,
		}}})
		guesses = append(guesses, models.PlayerGeoData{PlayerID: id, Rounds: geo})
	}
	geodata := &models.GameGeoData{PlayerGuesses: guesses}
	for r := 0; r < 5; r++ {
		geodata.ActualLocations = append(geodata.ActualLocations, models.RoundGeoData{
			Location: models.GeoData{Address: models.Address{Country: "United Arab Emirates", State: "Sharjah Emirate", CountryCode: "ae"}},
		})
	}

	embed := GenerateGeoGuessrDailyChallengeEmbed(items, challenge, geodata)
	assertEmbedLimits(t, embed)
}

func assertEmbedLimits(t *testing.T, embed *discordgo.MessageEmbed) {
	t.Helper()
	if len(embed.Fields) > maxEmbedFields {
		t.Errorf("too many fields: %d", len(embed.Fields))
	}
	for _, f := range embed.Fields {
		if len(f.Name) > maxFieldNameLen {
			t.Errorf("field name %q exceeds %d chars (%d)", f.Name, maxFieldNameLen, len(f.Name))
		}
		if len(f.Value) > maxFieldValueLen {
			t.Errorf("field %q value exceeds %d chars (%d)", f.Name, maxFieldValueLen, len(f.Value))
		}
	}
	if total := embedLength(embed); total > maxEmbedTotalLen {
		t.Errorf("embed total length %d exceeds %d", total, maxEmbedTotalLen)
	}
}
