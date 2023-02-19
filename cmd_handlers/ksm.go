package cmd_handlers

import "github.com/bwmarrin/discordgo"

func KsmHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Kim Sung-mo",
		},
	})
	if err != nil {
		return
	}
}
