package discord

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Miconen/geoguessr-daily-challenge-recap/models"
	"github.com/bwmarrin/discordgo"
)

type RoundPerformance struct {
	Round    int
	Score    int
	Time     int
	Distance float64
	Player   string
}

func createPlayerMap(players []models.ClubPlayer) map[string]models.ClubPlayer {
	m := make(map[string]models.ClubPlayer)
	for _, player := range players {
		m[player.ID] = player
	}
	return m
}

// GenerateGeoGuessrDailyChallengeEmbed builds a rich Discord embed with leaderboard, round-by-round recap, and highlights.
func GenerateGeoGuessrDailyChallengeEmbed(items []models.Items, challenge models.Challenge, geodata *models.GameGeoData) *discordgo.MessageEmbed {
	if len(items) == 0 {
		return &discordgo.MessageEmbed{
			Title:       "🌍 GeoGuessr Daily Challenge Recap",
			Description: "No player data available.",
			Color:       0x5865F2,
		}
	}

	playerLookup := createPlayerMap(challenge.Club)

	// Sort players by total score descending, then by total time ascending as tiebreaker
	sort.Slice(items, func(i, j int) bool {
		scoreI := parseScore(items[i].Game.Player.TotalScore.Amount)
		scoreJ := parseScore(items[j].Game.Player.TotalScore.Amount)
		if scoreI != scoreJ {
			return scoreI > scoreJ
		}
		return items[i].Game.Player.TotalTime < items[j].Game.Player.TotalTime
	})

	var dateStr string
	if !challenge.Date.IsZero() {
		dateStr = challenge.Date.Format("Jan 02, 2006")
	} else {
		dateStr = time.Now().Format("Jan 02, 2006")
	}

	embed := &discordgo.MessageEmbed{
		Title: fmt.Sprintf("🌍 GeoGuessr Daily Challenge Recap — %s", dateStr),
		Color: 0xFEE75C, // Gold color for competition recap
		Footer: &discordgo.MessageEmbedFooter{
			Text: "GeoGuessr Daily Challenge Bot",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// 1. Leaderboard Field
	leaderboard := buildLeaderboard(items, playerLookup)
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "🏆 Leaderboard",
		Value:  leaderboard,
		Inline: false,
	})

	// 2. Round-by-Round Breakdown Field
	roundsBreakdown := buildRoundByRoundSection(items, geodata)
	if roundsBreakdown != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "📍 Round-by-Round Breakdown",
			Value:  roundsBreakdown,
			Inline: false,
		})
	}

	// 3. Highlights Field
	highlights := buildPerformanceComparison(items)
	if highlights != "" {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "⚡ Highlights",
			Value:  highlights,
			Inline: false,
		})
	}

	return embed
}

// buildLeaderboard creates the overall standing overview for all players
func buildLeaderboard(items []models.Items, playerLookup map[string]models.ClubPlayer) string {
	var lines []string
	for i, item := range items {
		p := item.Game.Player
		cp := playerLookup[p.ID]

		streak := ""
		if cp.CurrentStreak > 0 {
			streak = fmt.Sprintf(" • 🔥 %d", cp.CurrentStreak)
		}

		flag := getFlagEmoji(p.CountryCode)
		totalScore := parseScore(p.TotalScore.Amount)
		distMeters := p.TotalDistanceInMeters
		if distMeters == 0 {
			distMeters = parseDistance(p.TotalDistance.Meters.Amount)
		}

		line := fmt.Sprintf("%s **%s** %s `%s pts` • `%s` • `%s`%s",
			getMedalEmoji(i),
			p.Nick,
			flag,
			formatNumber(totalScore),
			formatDistance(distMeters),
			formatTime(p.TotalTime),
			streak,
		)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

type PlayerRoundStat struct {
	PlayerName string
	Score      int
	Time       int
	Distance   float64
	GuessGeo   models.GeoData
}

// buildRoundByRoundSection groups players under each round
func buildRoundByRoundSection(items []models.Items, geodata *models.GameGeoData) string {
	if len(items) == 0 || geodata == nil {
		return ""
	}

	// Map player guesses by player ID
	playerGuessMap := make(map[string][]models.GuessGeoData)
	for _, pgd := range geodata.PlayerGuesses {
		playerGuessMap[pgd.PlayerID] = pgd.Rounds
	}

	numRounds := len(geodata.ActualLocations)
	if numRounds == 0 && len(items[0].Game.Rounds) > 0 {
		numRounds = len(items[0].Game.Rounds)
	}

	var sections []string

	for r := 0; r < numRounds; r++ {
		var actualLoc models.GeoData
		if r < len(geodata.ActualLocations) {
			actualLoc = geodata.ActualLocations[r].Location
		}

		locTitle := formatLocationTitle(actualLoc)
		roundHeader := fmt.Sprintf("**R%d:** %s", r+1, locTitle)

		// Collect player stats for this round
		var roundStats []PlayerRoundStat
		for _, item := range items {
			p := item.Game.Player
			if r >= len(p.Guesses) {
				continue
			}
			g := p.Guesses[r]
			dist := g.DistanceInMeters
			if dist == 0 {
				dist = parseDistance(g.Distance.Meters.Amount)
			}
			score := g.RoundScoreInPoints
			if score == 0 {
				score = parseScore(g.RoundScore.Amount)
			}

			var guessGeo models.GeoData
			if pRounds, ok := playerGuessMap[p.ID]; ok && r < len(pRounds) {
				guessGeo = pRounds[r].Guess
			}

			roundStats = append(roundStats, PlayerRoundStat{
				PlayerName: p.Nick,
				Score:      score,
				Time:       g.Time,
				Distance:   dist,
				GuessGeo:   guessGeo,
			})
		}

		// Sort round stats descending by score, ascending by time
		sort.Slice(roundStats, func(i, j int) bool {
			if roundStats[i].Score != roundStats[j].Score {
				return roundStats[i].Score > roundStats[j].Score
			}
			return roundStats[i].Time < roundStats[j].Time
		})

		var playerLines []string
		for idx, stat := range roundStats {
			branch := "├"
			if idx == len(roundStats)-1 {
				branch = "└"
			}

			var medal string
			if stat.Score == 5000 {
				medal = "⭐"
			} else {
				medal = getMedalEmoji(idx)
			}

			// Check wrong country guess
			wrongGuess := ""
			actCountryCode := strings.ToLower(actualLoc.Address.CountryCode)
			guessCountryCode := strings.ToLower(stat.GuessGeo.Address.CountryCode)
			if guessCountryCode != "" && actCountryCode != "" && guessCountryCode != actCountryCode {
				countryName := stat.GuessGeo.Address.Country
				if countryName == "" {
					countryName = strings.ToUpper(guessCountryCode)
				}
				wrongGuess = fmt.Sprintf(" ⤏ %s *%s*", getFlagEmoji(guessCountryCode), countryName)
			}

			line := fmt.Sprintf("%s %s **%s:** `%s pts` (%s • %s)%s",
				branch,
				medal,
				stat.PlayerName,
				formatNumber(stat.Score),
				formatDistance(stat.Distance),
				formatTime(stat.Time),
				wrongGuess,
			)
			playerLines = append(playerLines, line)
		}

		sections = append(sections, roundHeader+"\n"+strings.Join(playerLines, "\n"))
	}

	return strings.Join(sections, "\n\n")
}

func buildPerformanceComparison(items []models.Items) string {
	if len(items) == 0 {
		return ""
	}

	var (
		bestScore, worstScore RoundPerformance
		bestTime, worstTime   RoundPerformance
		perfects              = make(map[string]int)
	)

	worstScore.Score = 5001
	bestScore.Score = -1
	worstTime.Time = -1
	bestTime.Time = 999999

	for _, item := range items {
		p := item.Game.Player
		for i, g := range p.Guesses {
			score := g.RoundScoreInPoints
			if score == 0 {
				score = parseScore(g.RoundScore.Amount)
			}
			dist := g.DistanceInMeters
			if dist == 0 {
				dist = parseDistance(g.Distance.Meters.Amount)
			}

			if score > bestScore.Score {
				bestScore = RoundPerformance{Round: i + 1, Score: score, Time: g.Time, Distance: dist, Player: p.Nick}
			}
			if score < worstScore.Score {
				worstScore = RoundPerformance{Round: i + 1, Score: score, Time: g.Time, Distance: dist, Player: p.Nick}
			}
			if g.Time > worstTime.Time {
				worstTime = RoundPerformance{Round: i + 1, Score: score, Time: g.Time, Distance: dist, Player: p.Nick}
			}
			if g.Time < bestTime.Time {
				bestTime = RoundPerformance{Round: i + 1, Score: score, Time: g.Time, Distance: dist, Player: p.Nick}
			}

			if score == 5000 {
				perfects[p.Nick]++
			}
		}
	}

	var lines []string
	if bestScore.Score >= 0 {
		lines = append(lines, fmt.Sprintf("🎯 **Best Score:** %s (R%d • `%s pts`)", bestScore.Player, bestScore.Round, formatNumber(bestScore.Score)))
	}
	if worstScore.Score <= 5000 {
		lines = append(lines, fmt.Sprintf("⚠️ **Worst Score:** %s (R%d • `%s pts`)", worstScore.Player, worstScore.Round, formatNumber(worstScore.Score)))
	}
	if bestTime.Time < 999999 {
		lines = append(lines, fmt.Sprintf("⚡ **Fastest Round:** %s (R%d • `%s`)", bestTime.Player, bestTime.Round, formatTime(bestTime.Time)))
	}
	if worstTime.Time >= 0 {
		lines = append(lines, fmt.Sprintf("🐢 **Slowest Round:** %s (R%d • `%s`)", worstTime.Player, worstTime.Round, formatTime(worstTime.Time)))
	}

	if len(perfects) > 0 {
		var perfList []string
		for name, count := range perfects {
			perfList = append(perfList, fmt.Sprintf("%s (%d)", name, count))
		}
		lines = append(lines, fmt.Sprintf("⭐ **5k Perfects:** %s", strings.Join(perfList, ", ")))
	}

	return strings.Join(lines, "\n")
}

// --- Helpers ---

func formatLocationTitle(geo models.GeoData) string {
	flag := getFlagEmoji(geo.Address.CountryCode)
	country := geo.Address.Country
	state := geo.Address.State
	if state == "" {
		state = geo.Address.Province
	}
	if state == "" {
		state = geo.Address.County
	}

	if country == "" && state == "" {
		return "Unknown Location"
	}
	if state != "" && country != "" {
		return fmt.Sprintf("%s %s, %s", flag, country, state)
	}
	if country != "" {
		return fmt.Sprintf("%s %s", flag, country)
	}
	return fmt.Sprintf("%s %s", flag, state)
}

func getFlagEmoji(code string) string {
	code = strings.ToLower(strings.TrimSpace(code))
	if code == "" {
		return "🏳️"
	}
	return fmt.Sprintf(":flag_%s:", code)
}

func formatDistance(meters float64) string {
	if meters <= 0 {
		return "0 m"
	}
	if meters < 1000 {
		return fmt.Sprintf("%.0f m", meters)
	}
	km := meters / 1000.0
	if km < 10 {
		return fmt.Sprintf("%.1f km", km)
	}
	return fmt.Sprintf("%s km", formatNumber(int(km+0.5)))
}

func formatTime(s int) string {
	if s < 60 {
		return fmt.Sprintf("%ds", s)
	}
	return fmt.Sprintf("%dm %02ds", s/60, s%60)
}

func getMedalEmoji(i int) string {
	emojis := []string{"🥇", "🥈", "🥉"}
	if i >= 0 && i < len(emojis) {
		return emojis[i]
	}
	return fmt.Sprintf("%d.", i+1)
}

func parseScore(s string) int {
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

func parseDistance(s string) float64 {
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func formatNumber(n int) string {
	in := strconv.Itoa(n)
	var out []byte
	l := len(in)
	for i, c := range in {
		if i > 0 && (l-i)%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, byte(c))
	}
	return string(out)
}
