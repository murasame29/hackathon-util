package discord

import "github.com/bwmarrin/discordgo"

func (dh *DiscordHandler) Health(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "health ok!",
		},
	})
}
