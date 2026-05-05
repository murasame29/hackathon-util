package create

import "github.com/bwmarrin/discordgo"

const muteRoleDeny = discordgo.PermissionSendMessages |
	discordgo.PermissionCreatePublicThreads |
	discordgo.PermissionCreatePrivateThreads |
	discordgo.PermissionAddReactions |
	discordgo.PermissionVoiceConnect |
	discordgo.PermissionVoiceSpeak

func buildVCPermissionOverwrites(participantsRoleID, mentorRoleID, guildID string) []*discordgo.PermissionOverwrite {
	return []*discordgo.PermissionOverwrite{
		{ID: guildID, Type: discordgo.PermissionOverwriteTypeRole, Deny: discordgo.PermissionViewChannel},
		{ID: participantsRoleID, Type: discordgo.PermissionOverwriteTypeRole, Allow: discordgo.PermissionViewChannel},
		{ID: mentorRoleID, Type: discordgo.PermissionOverwriteTypeRole, Allow: discordgo.PermissionViewChannel},
	}
}

func buildCategoryPermissionOverwrites(teamRoleID, mentorRoleID, guildID, muteRoleID string) []*discordgo.PermissionOverwrite {
	overwrites := []*discordgo.PermissionOverwrite{
		{ID: guildID, Type: discordgo.PermissionOverwriteTypeRole, Deny: discordgo.PermissionViewChannel},
		{ID: teamRoleID, Type: discordgo.PermissionOverwriteTypeRole, Allow: discordgo.PermissionViewChannel},
		{ID: mentorRoleID, Type: discordgo.PermissionOverwriteTypeRole, Allow: discordgo.PermissionViewChannel},
	}
	if muteRoleID != "" {
		overwrites = append(overwrites, &discordgo.PermissionOverwrite{
			ID:    muteRoleID,
			Type:  discordgo.PermissionOverwriteTypeRole,
			Deny:  muteRoleDeny,
			Allow: 0,
		})
	}
	return overwrites
}

func buildPublicPermissionOverwrites(guildID string) []*discordgo.PermissionOverwrite {
	return []*discordgo.PermissionOverwrite{
		{ID: guildID, Type: discordgo.PermissionOverwriteTypeRole, Allow: discordgo.PermissionViewChannel},
	}
}
