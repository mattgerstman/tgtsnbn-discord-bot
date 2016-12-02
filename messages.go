package main

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mattgerstman/discordgo"
)

type GuildRoles map[string]string

var rolesMap = make(map[string]GuildRoles)

// Listener when a message is sent on discord.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.WithFields(log.Fields{
		"content":  m.Content,
		"username": m.Author.Username,
	}).Info("Message received.")

	if len(m.Mentions) == 0 {
		log.Info("No mentions, return early")
		return
	}
	if !strings.Contains(m.Content, "++") {
		log.Info("Didn't plus plus, return early")
		return
	}

	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Error("Unable to fetch channel from discord. ", err)
		return
	}
	guildID := channel.GuildID
	giver, err := s.GuildMember(guildID, m.Author.ID)
	if err != nil {
		log.Error("Unable to fetch giver from discord. ", err)
	}
	giverHouse, err := getHouseForMember(s, giver, guildID)
	if err != nil {
		log.Error("Unable to fetch house for giver from discord. ", err)
	}

	receiver, err := s.GuildMember(guildID, m.Mentions[0].ID)
	if err != nil {
		log.Error("Unable to fetch receiver from discord. ", err)
	}

	receiverHouse, err := getHouseForMember(s, receiver, guildID)
	if err != nil {
		log.Error("Unable to fetch house for receiver from discord. ", err)
	}

	if !canWeGivePoints(
		giver.User.ID,
		giverHouse,
		receiver.User.ID,
		receiverHouse,
	) {
		log.Info("Cannot give points return early")
		return
	}

	addPoints(m.Mentions[0].ID, guildID, receiverHouse)

	s.ChannelMessageSend(m.ChannelID, "10 Points to Slytherin")
}
