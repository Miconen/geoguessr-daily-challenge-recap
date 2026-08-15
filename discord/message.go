package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// SendDM sends a Discord embed message to a list of users
func SendDM(botToken string, users []string, embed *discordgo.MessageEmbed) error {
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
			fmt.Printf("error creating DM channel for user %s: %v\n", user, err)
			continue // Skip this user and continue with others
		}

		_, err = dg.ChannelMessageSendEmbed(channel.ID, embed)
		if err != nil {
			fmt.Printf("Error sending DM to user %s: %v\n", user, err)
			continue
		}
	}

	return nil
}
