package models

import "time"

type ClubPlayer struct {
	ID              string  `json:"id"`
	Nick            string  `json:"nick"`
	PinURL          string  `json:"pinUrl"`
	TotalScore      int     `json:"totalScore"`
	TotalTime       int     `json:"totalTime"`
	TotalDistance   float64 `json:"totalDistance"`
	IsOnLeaderboard bool    `json:"isOnLeaderboard"`
	IsVerified      bool    `json:"isVerified"`
	Flair           int     `json:"flair"`
	CountryCode     string  `json:"countryCode"`
	CurrentStreak   int     `json:"currentStreak"`
	TotalStepsCount int     `json:"totalStepsCount"`
}

type Challenge struct {
	Date         time.Time    `json:"date"`
	Description  string       `json:"description"`
	Participants int          `json:"participants"`
	Token        string       `json:"token"`
	Club         []ClubPlayer `json:"club"`
}

type Competition struct {
	Items           []Items `json:"items"`
	PaginationToken string  `json:"paginationToken"`
}

type Min struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Max struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Bounds struct {
	Min Min `json:"min"`
	Max Max `json:"max"`
}

type Rounds struct {
	Lat                float64   `json:"lat"`
	Lng                float64   `json:"lng"`
	PanoID             string    `json:"panoId"`
	Heading            float64   `json:"heading"`
	Pitch              float64   `json:"pitch"`
	Zoom               int       `json:"zoom"`
	StreakLocationCode string    `json:"streakLocationCode"`
	StartTime          time.Time `json:"startTime"`
}

type TotalScore struct {
	Amount     string  `json:"amount"`
	Unit       string  `json:"unit"`
	Percentage float64 `json:"percentage"`
}

type Meters struct {
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
}

type Miles struct {
	Amount string `json:"amount"`
	Unit   string `json:"unit"`
}

type TotalDistance struct {
	Meters Meters `json:"meters"`
	Miles  Miles  `json:"miles"`
}

type RoundScore struct {
	Amount     string  `json:"amount"`
	Unit       string  `json:"unit"`
	Percentage float64 `json:"percentage"`
}

type Distance struct {
	Meters Meters `json:"meters"`
	Miles  Miles  `json:"miles"`
}

type Guesses struct {
	Lat                    float64    `json:"lat"`
	Lng                    float64    `json:"lng"`
	TimedOut               bool       `json:"timedOut"`
	TimedOutWithGuess      bool       `json:"timedOutWithGuess"`
	SkippedRound           bool       `json:"skippedRound"`
	RoundScore             RoundScore `json:"roundScore"`
	RoundScoreInPercentage float64    `json:"roundScoreInPercentage"`
	RoundScoreInPoints     int        `json:"roundScoreInPoints"`
	Distance               Distance   `json:"distance"`
	DistanceInMeters       float64    `json:"distanceInMeters"`
	StepsCount             int        `json:"stepsCount"`
	StreakLocationCode     string     `json:"streakLocationCode"`
	Time                   int        `json:"time"`
}

type Pin struct {
	URL       string `json:"url"`
	Anchor    string `json:"anchor"`
	IsDefault bool   `json:"isDefault"`
}

type Player struct {
	TotalScore            TotalScore    `json:"totalScore"`
	TotalDistance         TotalDistance `json:"totalDistance"`
	TotalDistanceInMeters float64       `json:"totalDistanceInMeters"`
	TotalStepsCount       int           `json:"totalStepsCount"`
	TotalTime             int           `json:"totalTime"`
	TotalStreak           int           `json:"totalStreak"`
	Guesses               []Guesses     `json:"guesses"`
	IsLeader              bool          `json:"isLeader"`
	CurrentPosition       int           `json:"currentPosition"`
	Pin                   Pin           `json:"pin"`
	NewBadges             []string      `json:"newBadges"`
	Explorer              string        `json:"explorer"`
	ID                    string        `json:"id"`
	Nick                  string        `json:"nick"`
	IsVerified            bool          `json:"isVerified"`
	Flair                 int           `json:"flair"`
	CountryCode           string        `json:"countryCode"`
}

type CurrentLevel struct {
	Level   int `json:"level"`
	XpStart int `json:"xpStart"`
}

type NextLevel struct {
	Level   int `json:"level"`
	XpStart int `json:"xpStart"`
}

type CurrentTitle struct {
	ID           int    `json:"id"`
	TierID       int    `json:"tierId"`
	MinimumLevel int    `json:"minimumLevel"`
	Name         string `json:"name"`
}

type XpProgressions struct {
	Xp           int          `json:"xp"`
	CurrentLevel CurrentLevel `json:"currentLevel"`
	NextLevel    NextLevel    `json:"nextLevel"`
	CurrentTitle CurrentTitle `json:"currentTitle"`
}

type XpAwards struct {
	Xp     int    `json:"xp"`
	Reason string `json:"reason"`
	Count  int    `json:"count"`
}

type AwardedXp struct {
	TotalAwardedXp int        `json:"totalAwardedXp"`
	XpAwards       []XpAwards `json:"xpAwards"`
}

type ProgressChange struct {
	XpProgressions          []XpProgressions `json:"xpProgressions"`
	AwardedXp               AwardedXp        `json:"awardedXp"`
	Medal                   int              `json:"medal"`
	CompetitiveProgress     string           `json:"competitiveProgress"`
	RankedSystemProgress    string           `json:"rankedSystemProgress"`
	RankedTeamDuelsProgress string           `json:"rankedTeamDuelsProgress"`
	QuickplayDuelsProgress  string           `json:"quickplayDuelsProgress"`
}

type Game struct {
	Token            string         `json:"token"`
	Type             string         `json:"type"`
	Mode             string         `json:"mode"`
	State            string         `json:"state"`
	RoundCount       int            `json:"roundCount"`
	TimeLimit        int            `json:"timeLimit"`
	ForbidMoving     bool           `json:"forbidMoving"`
	ForbidZooming    bool           `json:"forbidZooming"`
	ForbidRotating   bool           `json:"forbidRotating"`
	GuessMapType     string         `json:"guessMapType"`
	StreakType       string         `json:"streakType"`
	Map              string         `json:"map"`
	MapName          string         `json:"mapName"`
	PanoramaProvider int            `json:"panoramaProvider"`
	Bounds           Bounds         `json:"bounds"`
	Round            int            `json:"round"`
	Rounds           []Rounds       `json:"rounds"`
	Player           Player         `json:"player"`
	ProgressChange   ProgressChange `json:"progressChange"`
}

type Items struct {
	Game Game `json:"game"`
}

// ClubLeaderboard is the response of challenges/daily-challenges/leaderboard/club?dateStr=YYYY-MM-DD.
type ClubLeaderboard struct {
	Token        string                 `json:"token"`
	TotalEntries int                    `json:"totalEntries"`
	Entries      []ClubLeaderboardEntry `json:"entries"`
}

type ClubLeaderboardEntry struct {
	Rank            int     `json:"rank"`
	UserID          string  `json:"userId"`
	Nick            string  `json:"nick"`
	PinURL          string  `json:"pinUrl"`
	CountryCode     string  `json:"countryCode"`
	Flair           int     `json:"flair"`
	CurrentStreak   int     `json:"currentStreak"`
	TotalScore      int     `json:"totalScore"`
	TotalDistance   float64 `json:"totalDistance"`
	TotalTime       int     `json:"totalTime"`
	TotalStepsCount int     `json:"totalStepsCount"`
	Medal           string  `json:"medal"`
	Game            Game    `json:"game"`
}

// ToChallenge converts the dated club leaderboard into the Challenge + Items shape used by the recap.
func (lb ClubLeaderboard) ToChallenge(date time.Time) (Challenge, []Items) {
	challenge := Challenge{Date: date, Token: lb.Token}
	items := make([]Items, 0, len(lb.Entries))
	for _, e := range lb.Entries {
		challenge.Club = append(challenge.Club, ClubPlayer{
			ID:              e.UserID,
			Nick:            e.Nick,
			PinURL:          e.PinURL,
			TotalScore:      e.TotalScore,
			TotalTime:       e.TotalTime,
			TotalDistance:   e.TotalDistance,
			Flair:           e.Flair,
			CountryCode:     e.CountryCode,
			CurrentStreak:   e.CurrentStreak,
			TotalStepsCount: e.TotalStepsCount,
		})
		items = append(items, Items{Game: e.Game})
	}
	return challenge, items
}
