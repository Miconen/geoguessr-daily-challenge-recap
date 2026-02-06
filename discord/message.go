package discord

import (
	"fmt"

	"github.com/Miconen/geoguessr-daily-challenge-recap/models"
	"github.com/bwmarrin/discordgo"
)

func SendDM(botToken string, users []string, challenge models.Challenge) error {
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return fmt.Errorf("error creating Discord session: %w", err)
	}

	err = dg.Open()
	if err != nil {
		return fmt.Errorf("error opening connection: %w", err)
	}
	defer dg.Close()

	// Send to each user
	for _, user := range users {
		channel, err := dg.UserChannelCreate(user)
		if err != nil {
			fmt.Printf("error creating DM channel for user %s: %v", user, err)
			continue // Skip this user and continue with others
		}

		embed := CreateChallengeEmbed(challenge)
		_, err = dg.ChannelMessageSendEmbed(channel.ID, embed)
		if err != nil {
			fmt.Printf("Error sending DM to user %s: %v", user, err)
			continue
		}
	}

	return nil
}
