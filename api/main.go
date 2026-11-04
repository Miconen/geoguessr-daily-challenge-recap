package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/Miconen/geoguessr-daily-challenge-recap/models"
)

var Endpoint = "https://www.geoguessr.com/api/v3/"

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

// FetchDailyChallenge fetches the club's results for the given day from the dated club leaderboard.
func FetchDailyChallenge(ncfa string, date time.Time) (models.Challenge, []models.Items, error) {
	lb, err := GeoGuessrRequest[models.ClubLeaderboard](ncfa, GetClubLeaderboardEndpoint(date))
	if err != nil {
		return models.Challenge{}, nil, fmt.Errorf("fetching club leaderboard for %s: %w", date.Format("2006-01-02"), err)
	}
	challenge, items := lb.ToChallenge(date)
	return challenge, items, nil
}
