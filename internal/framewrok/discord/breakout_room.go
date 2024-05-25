package discord

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/hackathon-util/cmd/config"
	// "github.com/murasame29/hackathon-util/internal/application"
)

var integerOptionMinValue = 1.0
var BreakoutRoomCommand = &discordgo.ApplicationCommand{
	Name:        "breakout-room",
	Description: "start breakout room",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionChannel,
			Name:        "target-vc",
			Description: "Host VC channel name",
			ChannelTypes: []discordgo.ChannelType{
				discordgo.ChannelTypeGuildVoice,
				discordgo.ChannelTypeGuildStageVoice,
			},
			Required: true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "number-per-room",
			Description: "Number of people in each room",
			MinValue:    &integerOptionMinValue,
			MaxValue:    100,
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "time-to-live",
			Description: "Breakout room time-to-live",
			MaxValue:    120,
			Required:    true,
		},
	},
}

func (dh *DiscordHandler) BreakoutRoom(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	targetVC := optionMap["target-vc"].ChannelValue(s)
	numberPerRoom := int(optionMap["number-per-room"].IntValue())
	timeToLive := int(optionMap["time-to-live"].IntValue())

	// Get the list of users in the target voice channel
	users, err := getUsersInChannel(s, targetVC.ID)
	if err != nil {
		log.Printf("Error getting users in channel: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to get users in the target channel.",
			},
		})
		return
	}

	// Calculate the number of breakout rooms needed
	numUsers := len(users)
	numRooms := (numUsers + numberPerRoom - 1) / numberPerRoom // Ceiling division
	log.Printf("Creating %d breakout rooms for %d users", numRooms, numUsers)
	breakoutRooms := make([]*discordgo.Channel, numRooms)

	// Create breakout rooms
	for j := 0; j < numRooms; j++ {
		room, err := s.GuildChannelCreate(targetVC.GuildID, fmt.Sprintf("breakout_%d", j+1), discordgo.ChannelTypeGuildVoice)
		if err != nil {
			log.Printf("Error creating breakout room: %v", err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to create breakout rooms.",
				},
			})
			return
		}
		breakoutRooms[j] = room
	}

	// Shuffle users and distribute them to breakout rooms
	rand.Shuffle(numUsers, func(i, j int) { users[i], users[j] = users[j], users[i] })
	for idx, user := range users {
		roomIdx := idx % numRooms
		err := s.GuildMemberMove(targetVC.GuildID, user, &breakoutRooms[roomIdx].ID)
		if err != nil {
			log.Printf("Error moving user to breakout room: %v", err)
		}
	}

	// Respond to the interaction
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Breakout rooms created and users assigned. They will be moved back in %d minutes.", timeToLive),
		},
	})

	// Wait for the time-to-live duration
	time.Sleep(time.Duration(timeToLive) * time.Minute)

	// Move users back to the original channel and delete breakout rooms
	for _, room := range breakoutRooms {
		users, err := getUsersInChannel(s, room.ID)
		if err != nil {
			log.Printf("Error getting users in breakout room: %v", err)
		}
		for _, user := range users {
			s.GuildMemberMove(targetVC.GuildID, user, &targetVC.ID)
		}
		// channel, err := s.ChannelDelete(room.ID)
		if err != nil {
			log.Printf("Error deleting breakout room: %v", err)
		}
		// log.Printf("Deleted breakout room: %v", channel.Name)
	}
}

func getUsersInChannel(s *discordgo.Session, channelID string) ([]string, error) {
	// Fetch the voice states for the guild
	guild, err := s.State.Guild(config.Config.Discord.GuildID[0])
	if err != nil {
		return nil, fmt.Errorf("failed to fetch guild: %v", err)
	}

	var users []string
	for _, vs := range guild.VoiceStates {
		if vs.ChannelID == channelID {
			users = append(users, vs.UserID)
		}
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("no users found in the channel")
	}
	return users, nil
}
