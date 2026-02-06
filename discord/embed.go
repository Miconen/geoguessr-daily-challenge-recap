package discord

import (
	"fmt"
	"math"
	"strings"

	"github.com/Miconen/geoguessr-daily-challenge-recap/models"
	"github.com/bwmarrin/discordgo"
)

func CreateChallengeEmbed(challenge models.Challenge) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:       "GeoGuessr Daily Results",
		URL:         fmt.Sprintf("https://www.geoguessr.com/results/%s", challenge.Token),
		Description: challenge.Date.Format("02.01.2006"),
		Color:       0x00ff00,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://play-lh.googleusercontent.com/DboQuoFNkqgfcl5NiLeXsSgUOLo1F_BMe0g9ZBQBFzq5GpX5M1o7LbJeMgocXmbfy8Y",
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("👥 %d participants", challenge.Participants),
		},
	}

	// Find best stats across all club players
	var bestScore int
	var bestDistance float64 = math.MaxFloat64
	var bestTime int = math.MaxInt

	for _, player := range challenge.Club {
		if player.TotalScore > bestScore {
			bestScore = player.TotalScore
		}
		if player.TotalDistance < bestDistance {
			bestDistance = player.TotalDistance
		}
		if player.TotalTime < bestTime {
			bestTime = player.TotalTime
		}
	}

	// Add club leaderboard with detailed stats
	if len(challenge.Club) > 0 {
		for i, player := range challenge.Club {
			if i >= 15 { // Reduced to avoid hitting embed limits
				break
			}

			embed.Fields = append(embed.Fields, createPlayerField(player, i+1, bestScore, bestDistance, bestTime))
		}
	} else {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Club Leaderboard",
			Value:  "No club players yet",
			Inline: false,
		})
	}

	return embed
}

func createPlayerField(player models.Player, position int, bestScore int, bestDistance float64, bestTime int) *discordgo.MessageEmbedField {
	// Position with medals
	positionStr := fmt.Sprintf("%d.", position)
	switch position {
	case 1:
		positionStr = "🥇"
	case 2:
		positionStr = "🥈"
	case 3:
		positionStr = "🥉"
	}

	// Build player name with badges
	var badges []string
	if player.CountryCode != "" {
		badges = append(badges, fmt.Sprintf(":flag_%s:", strings.ToLower(player.CountryCode)))
	}

	badgeStr := ""
	if len(badges) > 0 {
		badgeStr = " " + strings.Join(badges, " ")
	}

	name := fmt.Sprintf("%s %s%s", positionStr, player.Nick, badgeStr)

	// Build stats in a more readable format
	var stats []string

	// Secondary stats
	var secondaryStats []string

	if player.TotalStepsCount == 0 {
		secondaryStats = append(secondaryStats, "🚫 **No Move**")
	} else {
		secondaryStats = append(secondaryStats, fmt.Sprintf("👣 %d steps", player.TotalStepsCount))
	}

	if player.CurrentStreak > 0 {
		secondaryStats = append(secondaryStats, fmt.Sprintf("🔥 %d streak", player.CurrentStreak))
	}

	if len(secondaryStats) > 0 {
		stats = append(stats, strings.Join(secondaryStats, "  |  "))
	}

	// Primary stats on one line with best stat highlighting
	scoreStr := formatScore(player.TotalScore)
	if player.TotalScore == bestScore {
		scoreStr = fmt.Sprintf("**%s**", scoreStr)
	}

	distanceStr := formatDistance(player.TotalDistance)
	if player.TotalDistance == bestDistance {
		distanceStr = fmt.Sprintf("**%s**", distanceStr)
	}

	timeStr := formatTime(player.TotalTime)
	if player.TotalTime == bestTime {
		timeStr = fmt.Sprintf("**%s**", timeStr)
	}

	primaryStats := fmt.Sprintf("**Score:** %s  •  **Distance:** %s  •  **Time:** %s",
		scoreStr,
		distanceStr,
		timeStr)
	stats = append(stats, primaryStats)

	return &discordgo.MessageEmbedField{
		Name:   name,
		Value:  strings.Join(stats, "\n"),
		Inline: false,
	}
}

func formatScore(score int) string {
	if score >= 25000 {
		return fmt.Sprintf("%d 🎯", score)
	}
	return fmt.Sprintf("%d", score)
}

func formatDistance(distance float64) string {
	km := distance / 1000
	if km < 1 {
		return fmt.Sprintf("%.0f m", distance)
	} else if km < 10 {
		return fmt.Sprintf("%.2f km", km)
	} else if km < 100 {
		return fmt.Sprintf("%.1f km", km)
	}
	return fmt.Sprintf("%.0f km", km)
}

func formatTime(seconds int) string {
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := seconds / 60
	remainingSeconds := seconds % 60

	if minutes < 60 {
		if remainingSeconds == 0 {
			return fmt.Sprintf("%dm", minutes)
		}
		return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
	}

	hours := minutes / 60
	remainingMinutes := minutes % 60

	if remainingMinutes == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh %dm", hours, remainingMinutes)
}
