package main

import "github.com/mattgerstman/discordgo"

var rolesMap = make(map[string]GuildRoles)

/**
 * Gets the name of a role.
 */
func GetRoleName(
	s *discordgo.Session,
	roleId string,
	guildID string,
) (string, *ApplicationError) {
	if rolesMap[guildID] != nil {
		return rolesMap[guildID][roleId], nil
	}

	guildRoles, err := s.GuildRoles(guildID)
	if err != nil {
		return "",
			NewApplicationError(
				"Error getting list of roles for guild.",
				err,
				ErrorFetchRoles,
			)
	}

	rolesMap[guildID] = make(GuildRoles)
	for _, role := range guildRoles {
		rolesMap[guildID][role.ID] = role.Name
	}

	return rolesMap[guildID][roleId], nil
}
