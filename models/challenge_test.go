package models

import (
	"encoding/json"
	"testing"
	"time"
)

// Trimmed from a real challenges/daily-challenges/leaderboard/club response (2026-11-04).
const clubLeaderboardJSON = `{"endTime":"2026-11-04T23:59:59.0000000Z","token":"Av4U80pOdUfaj3yK","totalEntries":1,"entries":[
{"rank":1,"userId":"6066c99dfc6ef0000194aea4","nick":"Cerkie","countryCode":"dk","currentStreak":132,"totalScore":21913,"totalDistance":1065679.43,"totalTime":514,"medal":"Silver",
"game":{"token":"N1b6tKDk7nB6kBWu","roundCount":5,
"rounds":[{"lat":6.96,"lng":80.76,"streakLocationCode":"lk","region":{"territory":null,"country":"lk","continent":"as","hemisphere":"Northern"}}],
"player":{"totalScoreInPoints":21913,"totalDistanceInMeters":1065679.43,"totalTime":513,"id":"6066c99dfc6ef0000194aea4","nick":"Cerkie","countryCode":"dk",
"guesses":[{"lat":6.96,"lng":80.76,"roundScoreInPoints":5000,"distanceInMeters":3.5,"time":150,"streakLocationCode":null,
"regionGuessBonuses":[{"region":"Hemisphere","correct":true},{"region":"Continent","correct":true},{"region":"Country","correct":false}]}]}}}]}`

func TestDecodeClubLeaderboard(t *testing.T) {
	var lb ClubLeaderboard
	if err := json.Unmarshal([]byte(clubLeaderboardJSON), &lb); err != nil {
		t.Fatal(err)
	}
	challenge, items := lb.ToChallenge(time.Date(2026, 11, 4, 0, 0, 0, 0, time.UTC))
	if len(items) != 1 || len(challenge.Club) != 1 {
		t.Fatalf("got %d items, %d club players", len(items), len(challenge.Club))
	}
	game := items[0].Game
	if got := game.Player.TotalScoreInPoints; got != 21913 {
		t.Errorf("TotalScoreInPoints = %d, want 21913", got)
	}
	if got := challenge.Club[0].Rank; got != 1 {
		t.Errorf("Rank = %d, want 1", got)
	}
	if got := game.Rounds[0].Region.Country; got != "lk" {
		t.Errorf("Region.Country = %q, want lk", got)
	}
	correct, ok := game.Player.Guesses[0].CountryCorrect()
	if !ok || correct {
		t.Errorf("CountryCorrect() = %v, %v; want false, true", correct, ok)
	}
}

func TestCountryCorrectMissing(t *testing.T) {
	if _, ok := (Guesses{}).CountryCorrect(); ok {
		t.Error("expected ok=false without a Country bonus")
	}
}
