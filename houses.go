package main

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/mattgerstman/discordgo"
)

func getHouseForMember(s *discordgo.Session, member *discordgo.Member, guildID string) (string, error) {
	for _, role := range member.Roles {
		log.WithFields(log.Fields{
			"guildID": guildID,
			"role":    role,
		}).Info("Getting house for member")

		role, err := getRoleName(s, role, guildID)
		if err != nil {
			return "", err
		}

		houses := GetHouseMap()
		if _, ok := houses[role]; ok {
			return role, nil
		}
	}
	return "", errors.New("Role not found")
}
