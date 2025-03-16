package discord

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/hackathon-util/pkg/cache"
)

type Member struct {
	ss *discordgo.Session

	cache *cache.Cache
}

func newMember(ss *discordgo.Session) *Member {
	return &Member{ss: ss, cache: cache.New()}
}

func (m *Member) Get(ctx context.Context, guildID, name string) (*discordgo.Member, error) {
	if member, ok := m.cache.Get(fmt.Sprintf("%s_%s", guildID, name)); ok {
		return member.(*discordgo.Member), nil
	}

	var highestID int

	for {
		members, err := m.ss.GuildMembers(guildID, strconv.Itoa(highestID), 1000, discordgo.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		if len(members) == 0 {
			break
		}

		for _, member := range members {
			m.cache.Set(fmt.Sprintf("%s_%s", guildID, member.User.GlobalName), member)
			id, _ := strconv.Atoi(member.User.ID)
			if id > highestID {
				highestID = id
			}
		}
	}

	if member, ok := m.cache.Get(fmt.Sprintf("%s_%s", guildID, name)); ok {
		return member.(*discordgo.Member), nil
	}

	return nil, ErrResourceNotFound
}

func (m *Member) AddRole(ctx context.Context, guildID, userID, roleID string) error {
	return m.ss.GuildMemberRoleAdd(guildID, userID, roleID, discordgo.WithContext(ctx))
}
