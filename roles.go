package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mattgerstman/discordgo"
)

func getRoleName(s *discordgo.Session, roleId string, guildID string) (string, error) {
	if rolesMap[guildID] != nil {
		return rolesMap[guildID][roleId], nil
	}

	guildRoles, err := s.GuildRoles(guildID)
	if err != nil {
		log.Error(err)
		return "", err
	}

	rolesMap[guildID] = make(GuildRoles)
	for _, role := range guildRoles {
		rolesMap[guildID][role.ID] = role.Name
	}

	return rolesMap[guildID][roleId], nil
}
