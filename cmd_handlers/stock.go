package cmd_handlers

import "github.com/bwmarrin/discordgo"

func StockHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Stock",
		},
	})
	if err != nil {
		return
	}
}
