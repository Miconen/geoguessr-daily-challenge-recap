package models

import "time"

type Player struct {
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

type Author struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	AvatarImage       string   `json:"avatarImage"`
	CustomName        string   `json:"customName"`
	CustomAvatarImage string   `json:"customAvatarImage"`
	SignupAssetIds    []string `json:"signupAssetIds"`
	SignupCoins       int      `json:"signupCoins"`
	YoutubeLink       string   `json:"youtubeLink"`
	TwitchLink        string   `json:"twitchLink"`
	TwitterLink       string   `json:"twitterLink"`
	InstagramLink     string   `json:"instagramLink"`
	Program           string   `json:"program"`
}

type Challenge struct {
	AuthorCreator Author    `json:"authorCreator"`
	Date          time.Time `json:"date"`
	Description   string    `json:"description"`
	Participants  int       `json:"participants"`
	Token         string    `json:"token"`
	PickedWinner  bool      `json:"pickedWinner"`
	Leaderboard   []Player  `json:"leaderboard"`
	Friends       []Player  `json:"friends"`
	Country       []Player  `json:"country"`
	Club          []Player  `json:"club"`
}
