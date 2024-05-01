package discordgo

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/cockroachdb/errors"
	"github.com/murasame29/hackathon-util/cmd/config"

	"k8s.io/utils/ptr"
)

func (ds *DiscordSession) CreateRole(ctx context.Context, roleName string) error {
	_, err := ds.ss.GuildRoleCreate(config.Config.Discord.GuildID, &discordgo.RoleParams{
		Name:  roleName,
		Color: ptr.To(640000),
	}, discordgo.WithContext(ctx))

	if err != nil {
		return errors.Wrapf(err, "error creating role %s", roleName)
	}

	return nil
}

func (ds *DiscordSession) GetRoles(ctx context.Context) (map[string]string, error) {
	roles, err := ds.ss.GuildRoles(config.Config.Discord.GuildID, discordgo.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "error getting roles")
	}

	roleMap := make(map[string]string)
	for _, role := range roles {
		roleMap[role.Name] = role.ID
	}

	return roleMap, nil
}

func (ds *DiscordSession) DeleteRole(ctx context.Context, roleID string) error {
	err := ds.ss.GuildRoleDelete(config.Config.Discord.GuildID, roleID, discordgo.WithContext(ctx))

	if err != nil {
		return errors.Wrapf(err, "error deleting role %s", roleID)
	}

	return nil
}

func (ds *DiscordSession) BindRole(ctx context.Context, serverID, userID, roleID string) error {
	err := ds.ss.GuildMemberRoleAdd(serverID, userID, roleID, discordgo.WithContext(ctx))
	if err != nil {
		return err
	}

	return nil
}
