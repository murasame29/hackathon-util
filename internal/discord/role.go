package discord

import (
	"context"
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/murasame29/hackathon-util/pkg/cache"
	"github.com/murasame29/hackathon-util/pkg/ptr"
)

type Role struct {
	ss    *discordgo.Session
	cache *cache.Cache
}

func newRole(ss *discordgo.Session) *Role {
	return &Role{ss: ss, cache: cache.New()}
}

func (r *Role) Get(ctx context.Context, guildID, name string) (*discordgo.Role, error) {
	if channel, ok := r.cache.Get(fmt.Sprintf("%s_%s", guildID, name)); ok {
		return channel.(*discordgo.Role), nil
	}

	roles, err := r.ss.GuildRoles(guildID, discordgo.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		r.cache.Set(fmt.Sprintf("%s_%s", guildID, name), role)
		if role.Name == name {
			return role, nil
		}
	}

	return nil, ErrResourceNotFound
}

func (r *Role) Create(ctx context.Context, guildID, name string) (*discordgo.Role, error) {
	ok, err := r.Exist(ctx, guildID, name)
	if err != nil {
		return nil, err
	}

	if ok {
		return nil, ErrResourceAlreadyExists
	}

	result, err := r.ss.GuildRoleCreate(guildID, &discordgo.RoleParams{
		Name:        name,
		Mentionable: ptr.To(true),
	}, discordgo.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Role) Exist(ctx context.Context, guildID, name string) (bool, error) {
	_, err := r.Get(ctx, guildID, name)
	if err != nil {
		if errors.Is(err, ErrResourceNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
