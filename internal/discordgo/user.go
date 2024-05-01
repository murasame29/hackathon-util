package discordgo

import (
	"context"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/cockroachdb/errors"
)

func (ds *DiscordSession) GetUsers(ctx context.Context, serverID, after string) ([]*discordgo.Member, error) {
	users, err := ds.ss.GuildMembers(serverID, after, 1000, discordgo.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get users")
	}

	return users, nil
}

func (ds *DiscordSession) GetUsersAll(ctx context.Context, serverID string) (map[string]string, error) {
	users, err := ds.GetUsers(ctx, serverID, "")
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]string)
	highestID := 0

	highestID, userMap = highestUserID(users, userMap)

	for {
		users, err := ds.GetUsers(ctx, serverID, strconv.Itoa(highestID))
		if err != nil {
			return nil, err
		}

		if len(users) == 0 {
			break
		}

		highestID, userMap = highestUserID(users, userMap)
	}

	return userMap, nil
}

func highestUserID(members []*discordgo.Member, userMap map[string]string) (int, map[string]string) {
	var highestUserID int
	for _, user := range members {
		userMap[user.User.Username] = user.User.ID
		id, _ := strconv.Atoi(user.User.ID)
		if id > highestUserID {
			highestUserID = id
		}
	}

	return highestUserID, userMap
}
