package api

import (
	"testing"
	"time"
)

func TestParseChallengeDate(t *testing.T) {
	now := time.Date(2026, 9, 8, 15, 30, 0, 0, time.FixedZone("EEST", 3*3600))
	cases := map[string]string{
		"":           "2026-09-08",
		"today":      "2026-09-08",
		"yesterday":  "2026-09-07",
		"2026-09-06": "2026-09-06",
	}
	for in, want := range cases {
		got, err := ParseChallengeDate(in, now)
		if err != nil {
			t.Fatalf("%q: unexpected error %v", in, err)
		}
		if got.Format("2006-01-02") != want {
			t.Errorf("%q: got %s want %s", in, got.Format("2006-01-02"), want)
		}
	}
	if _, err := ParseChallengeDate("06/09/2026", now); err == nil {
		t.Error("expected error for bad date format")
	}
}

func TestGetClubLeaderboardEndpoint(t *testing.T) {
	d := time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC)
	want := Endpoint + "challenges/daily-challenges/leaderboard/club?dateStr=2026-09-06"
	if got := GetClubLeaderboardEndpoint(d); got != want {
		t.Errorf("got %s want %s", got, want)
	}
}
