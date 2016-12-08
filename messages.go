package main

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mattgerstman/discordgo"
)

type GuildRoles map[string]string

/**
 * Add points to a user's account.
 */
func plusPlus(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	giverHouse, appErr := GetHouseForMember(s, giver, guildID)
	if appErr != nil {
		log.Error("Unable to fetch house for giver from discord. ", appErr)
	}

	receiver, err := s.GuildMember(guildID, m.Mentions[0].ID)
	if err != nil {
		log.Error("Unable to fetch receiver from discord. ", err)
	}

	receiverID := receiver.User.ID
	if userIsThisBot(receiverID) {
		log.Info("User giving points to dumbledore")
		s.ChannelMessageSend(m.ChannelID, "Thank you!")
		return
	}

	receiverHouse, appErr := GetHouseForMember(s, receiver, guildID)
	if appErr != nil {
		log.Error("Unable to fetch house for receiver from discord. ", appErr)
	}

	giverID := giver.User.ID
	if giverID == receiverID {
		log.Info("Receiver cannot be giver")
		s.ChannelMessageSend(m.ChannelID, "Ten points to Dumbledore!")
		return
	}

	if giverHouse == receiverHouse {
		log.Info("Cannot give points to members of their own house")
		message := fmt.Sprintf(
			"A %s would give points to another %s. Fifty points to Buckbeak.",
			giverHouse,
			giverHouse,
		)
		s.ChannelMessageSend(m.ChannelID, message)
		return
	}

	appErr = AddPoints(m.Mentions[0].ID, guildID, receiverHouse)
	if appErr != nil {
		log.Error("Error adding points ", appErr)
	}

	userPoints, appErr := GetPointsForUser(receiverID, guildID)
	if appErr != nil {
		log.Error("Error getting user points ", appErr)
	}

	housePoints, appErr := GetPointsForHouse(receiverHouse, guildID)
	if appErr != nil {
		log.Error("Error getting house points ", appErr)
	}

	receiverName := GetNameForMember(receiver)
	response := fmt.Sprintf(
		"Ten points to %s in %s. %s now has %d points. %s now has %d points.",
		receiverName,
		receiverHouse,
		receiverName,
		userPoints,
		receiverHouse,
		housePoints,
	)
	s.ChannelMessageSend(m.ChannelID, response)
}

/**
 * Print the leaderboard.
 */
func leaderboard(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Error("Unable to fetch channel from discord. ", err)
		return
	}

	if !includesThisBot(m.Mentions) {
		log.Info("Didn't mention bot, return")
		return
	}

	housePoints, appErr := GetHouseLeaderboard(channel.GuildID)
	if appErr != nil {
		log.Error(appErr)
		return
	}

	message := "The points stand thus: "
	for _, housePoint := range housePoints {
		message += fmt.Sprintf("\n%s:\t%s", housePoint.House, housePoint.NumPoints)
	}

	s.ChannelMessageSend(m.ChannelID, message)
}

/**
 * Helper to determine if a list of users includes the running bot.
 */
func includesThisBot(users []*discordgo.User) bool {
	for _, mentionedUser := range users {
		if userIsThisBot(mentionedUser.ID) {
			return true
		}
	}
	return false
}

/**
 * Helper to determine if a single user is the running bot.
 */
func userIsThisBot(userID string) bool {
	config := GetConfig()
	return config.BotId == userID
}

/*
 * Listener when a message is sent on discord.
 */
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.WithFields(log.Fields{
		"content":  m.Content,
		"username": m.Author.Username,
	}).Info("Message received.")

	if m.Author.Bot {
		log.Info("Author is a bot, return early")
		return
	}

	if len(m.Mentions) == 0 {
		log.Info("No mentions, return early")
		return
	}

	if strings.Contains(m.Content, "++") {
		log.Info("Plus plus, add points")
		plusPlus(s, m)
		return
	}

	if strings.Contains(m.Content, "leaderboard") {
		log.Info("Leaderboard")
		leaderboard(s, m)
		return
	}

	log.Info("No function in message, don't do anything")
}
