package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mattgerstman/discordgo"
)

// Gets a member's roles and figure's out which one is their house.
func GetHouseForMember(
	s *discordgo.Session,
	member *discordgo.Member,
	guildID string,
) (string, *ApplicationError) {
	for _, role := range member.Roles {
		log.WithFields(log.Fields{
			"guildID": guildID,
			"role":    role,
		}).Info("Getting house for member")

		roleName, appErr := GetRoleName(s, role, guildID)
		if appErr != nil {
			return "", appErr
		}

		houses := GetHouseMap()
		if _, ok := houses[roleName]; ok {
			return roleName, nil
		}
	}
	return "", NewApplicationErrorWithoutError("Role not found", ErrorRoleNotFound)
}
