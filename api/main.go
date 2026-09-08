package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/Miconen/geoguessr-daily-challenge-recap/models"
)

var Endpoint = "https://www.geoguessr.com/api/v3/"
var EndpointDaily = Endpoint + "challenges/daily-challenges/today"

func GetScoresEndpoint(id string) string {
	return fmt.Sprintf(Endpoint+"results/highscores/%s?friends=false&limit=26&minRounds=5&club=true", id)
}

// GetClubLeaderboardEndpoint returns the club leaderboard (with full games) for a given UTC day.
func GetClubLeaderboardEndpoint(date time.Time) string {
	return Endpoint + "challenges/daily-challenges/leaderboard/club?dateStr=" + date.UTC().Format("2006-01-02")
}

// ParseChallengeDate accepts "today", "yesterday", or YYYY-MM-DD (UTC).
func ParseChallengeDate(s string, now time.Time) (time.Time, error) {
	day := now.UTC().Truncate(24 * time.Hour)
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "today":
		return day, nil
	case "yesterday":
		return day.AddDate(0, 0, -1), nil
	}
	d, err := time.Parse("2006-01-02", strings.TrimSpace(s))
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date %q: use today, yesterday or YYYY-MM-DD", s)
	}
	return d, nil
}

// FetchDailyChallenge fetches the club's results for the given day.
// Today uses the live "today" + highscores endpoints; any other day uses the dated club leaderboard.
func FetchDailyChallenge(ncfa string, date time.Time, now time.Time) (models.Challenge, []models.Items, error) {
	today := now.UTC().Truncate(24 * time.Hour)
	if date.Equal(today) {
		challenge, err := GeoGuessrRequest[models.Challenge](ncfa, EndpointDaily)
		if err != nil {
			return models.Challenge{}, nil, fmt.Errorf("fetching today's challenge: %w", err)
		}
		competition, err := GeoGuessrRequest[models.Competition](ncfa, GetScoresEndpoint(challenge.Token))
		if err != nil {
			return models.Challenge{}, nil, fmt.Errorf("fetching highscores: %w", err)
		}
		return challenge, competition.Items, nil
	}

	lb, err := GeoGuessrRequest[models.ClubLeaderboard](ncfa, GetClubLeaderboardEndpoint(date))
	if err != nil {
		return models.Challenge{}, nil, fmt.Errorf("fetching club leaderboard for %s: %w", date.Format("2006-01-02"), err)
	}
	challenge, items := lb.ToChallenge(date)
	return challenge, items, nil
}
