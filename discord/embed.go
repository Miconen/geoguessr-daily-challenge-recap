package discord

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Miconen/geoguessr-daily-challenge-recap/models"
)

// Simplified Data Structures

type RoundPerformance struct {
	Round    int
	Score    int
	Time     int
	Distance string
	Player   string
}

func createPlayerMap(players []models.ClubPlayer) map[string]models.ClubPlayer {
	m := make(map[string]models.ClubPlayer)
	for _, player := range players {
		m[player.ID] = player
	}
	return m
}

// func GenerateGeoGuessrDailyChallengeEmbed(items []models.Items, challenge models.Challenge, geodata *models.GameGeoData) *discordgo.MessageEmbed {
func GenerateGeoGuessrDailyChallengeEmbed(items []models.Items, challenge models.Challenge, geodata *models.GameGeoData) string {
	// Create the lookup map
	playerLookup := createPlayerMap(challenge.Club)

	// Sort by score descending
	sort.Slice(items, func(i, j int) bool {
		return items[i].Game.Player.TotalScore.Amount > items[j].Game.Player.TotalScore.Amount
	})

	bestScoresPerRound := findBestScoresPerRound(items)

	var msg []string

	// Player Breakdowns
	for i, item := range items {
		msg = append(msg, fmt.Sprintf("%s\n\n**Rounds:**\n%s", formatPlayerStats(i, item.Game.Player, playerLookup[item.Game.Player.ID]), buildRoundsBreakdown(item.Game.Player, bestScoresPerRound, geodata)))
	}

	// Performance Comparisons
	msg = append(msg, buildPerformanceComparison(items))

	return strings.Join(msg, "\n")
}

// --- Reusable Helpers ---

func formatPlayerStats(placement int, p models.Player, cp models.ClubPlayer) string {
	streak := ""
	if cp.CurrentStreak > 0 {
		streak = fmt.Sprintf("• 🔥 %d", cp.CurrentStreak)
	}
	return fmt.Sprintf("# %s **%s** :flag_%s: %s\n❯❯ %s pts • %s km • %s",
		getMedalEmoji(placement), p.Nick, p.CountryCode, streak, p.TotalScore.Amount, p.TotalDistance.Meters.Amount, formatTime(p.TotalTime))
}

func formatTime(s int) string {
	if s < 60 {
		return fmt.Sprintf("%ds", s)
	}
	return fmt.Sprintf("%dm %ds", s/60, s%60)
}

func getMedalEmoji(i int) string {
	emojis := []string{"🥇", "🥈", "🥉"}
	if i < 3 {
		return emojis[i]
	}
	return fmt.Sprintf("%d.", i+1)
}

func getNumberEmoji(n int) string {
	if n >= 1 && n <= 10 {
		return []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟"}[n-1]
	}
	return fmt.Sprintf("%d", n)
}

func findBestScoresPerRound(items []models.Items) []int {
	if len(items) == 0 {
		return nil
	}
	best := make([]int, len(items[0].Game.Player.Guesses))
	for _, item := range items {
		for i, g := range item.Game.Player.Guesses {
			if g.RoundScoreInPoints > best[i] {
				best[i] = g.RoundScoreInPoints
			}
		}
	}
	return best
}

func buildRoundsBreakdown(Player models.Player, bestScores []int, geodata *models.GameGeoData) string {
	var res string
	for i, g := range Player.Guesses {
		isBest := i < len(bestScores) && g.RoundScoreInPoints == bestScores[i] && g.RoundScoreInPoints > 0
		star := ""
		if g.RoundScoreInPoints == 5000 {
			star = "⭐ "
		}

		playerMap := make(map[string][]models.GuessGeoData)
		for _, pgd := range geodata.PlayerGuesses {
			playerMap[pgd.PlayerID] = pgd.Rounds
		}

		l := geodata.ActualLocations[i].Location
		c := playerMap[Player.ID][i].Guess

		loc := fmt.Sprintf(":flag_%s: %s, %s", l.Address.CountryCode, l.Address.Country, l.Address.State)
		if l.Address.CountryCode != c.Address.CountryCode {
			loc += fmt.Sprintf(" | :round_pushpin: :flag_%s: %s, %s", c.Address.CountryCode, c.Address.Country, c.Address.State)
		}

		line := fmt.Sprintf("%s%s %s\n", star, getNumberEmoji(i+1), loc)
		line += fmt.Sprintf("❯❯ %s pts • %s km • %s", g.RoundScore.Amount, g.Distance.Meters.Amount, formatTime(g.Time))

		if isBest {
			line = "**" + line + "**"
		}
		res += line + "\n"
	}
	return res
}

func buildPerformanceComparison(items []models.Items) string {
	if len(items) == 0 {
		return "No data"
	}

	var (
		bestScore, worstScore RoundPerformance
		bestTime, worstTime   RoundPerformance
		totalTime             int
		perfects              = make(map[string]int)
	)
	worstScore.Score = 5000
	worstTime.Time = 0
	bestTime.Time = 9999

	for _, item := range items {
		p := item.Game.Player
		totalTime += p.TotalTime
		if p.TotalTime < bestTime.Time {
			bestTime.Time, bestTime.Player = p.TotalTime, p.Nick
		}

		for i, g := range p.Guesses {
			if g.RoundScoreInPoints >= bestScore.Score {
				bestScore = RoundPerformance{
					Round:    i + 1,
					Score:    g.RoundScoreInPoints,
					Time:     g.Time,
					Distance: g.Distance.Meters.Amount,
					Player:   p.Nick,
				}
			}
			if g.RoundScoreInPoints <= worstScore.Score {
				worstScore = RoundPerformance{
					Round:    i + 1,
					Score:    g.RoundScoreInPoints,
					Time:     g.Time,
					Distance: g.Distance.Meters.Amount,
					Player:   p.Nick,
				}
			}
			if g.Time >= worstTime.Time {
				worstTime = RoundPerformance{
					Round:    i + 1,
					Score:    g.RoundScoreInPoints,
					Time:     g.Time,
					Distance: g.Distance.Meters.Amount,
					Player:   p.Nick,
				}
			}
			if g.Time <= bestTime.Time {
				bestTime = RoundPerformance{
					Round:    i + 1,
					Score:    g.RoundScoreInPoints,
					Time:     g.Time,
					Distance: g.Distance.Meters.Amount,
					Player:   p.Nick,
				}
			}

			if g.RoundScoreInPoints == 5000 {
				perfects[p.Nick]++
			}
		}
	}

	res := ""

	var scores []string
	scores = append(scores, fmt.Sprintf("**Best Score:** %s %s %d pts", bestScore.Player, getNumberEmoji(bestScore.Round), bestScore.Score))
	scores = append(scores, fmt.Sprintf("**Worst Score:** %s %s %d pts\n", worstScore.Player, getNumberEmoji(worstScore.Round), worstScore.Score))
	res += strings.Join(scores, "\n")

	var times []string
	times = append(times, fmt.Sprintf("**Fastest Time:** %s %s %s", bestTime.Player, getNumberEmoji(bestTime.Round), formatTime(bestTime.Time)))
	times = append(times, fmt.Sprintf("**Slowest Time:** %s %s %s\n", worstTime.Player, getNumberEmoji(worstTime.Round), formatTime(worstTime.Time)))
	res += strings.Join(times, "\n")

	res += "**Perfects:** "
	if len(perfects) == 0 {
		res += "None"
	}
	for name, count := range perfects {
		res += fmt.Sprintf("%s: %d ", name, count)
	}

	return res
}
