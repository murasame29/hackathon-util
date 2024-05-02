package discordgo

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/cockroachdb/errors"

	"k8s.io/utils/ptr"
)

func (ds *DiscordSession) CreateRole(ctx context.Context, serverID, roleName string) error {
	_, err := ds.ss.GuildRoleCreate(serverID, &discordgo.RoleParams{
		Name:  roleName,
		Color: ptr.To(640000),
	}, discordgo.WithContext(ctx))

	if err != nil {
		return errors.Wrapf(err, "error creating role %s", roleName)
	}

	return nil
}

func (ds *DiscordSession) GetRoles(ctx context.Context, serverID string) (map[string]string, error) {
	roles, err := ds.ss.GuildRoles(serverID, discordgo.WithContext(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "error getting roles")
	}

	roleMap := make(map[string]string)
	for _, role := range roles {
		roleMap[role.Name] = role.ID
	}

	return roleMap, nil
}

func (ds *DiscordSession) DeleteRole(ctx context.Context, serverID, roleID string) error {
	err := ds.ss.GuildRoleDelete(serverID, roleID, discordgo.WithContext(ctx))

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
