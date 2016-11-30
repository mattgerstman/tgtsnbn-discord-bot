package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/mattgerstman/discordgo"
)

func main() {
	log.Info("Hello Harry Potter")
	GetDB()
	config := GetConfig()

	discord, err := discordgo.New(config.Username, config.Password)
	if err != nil {
		log.Fatal("Failed to connect to discord", err)
	}

	header := http.Header{}
	header.Add("accept-encoding", "zlib")
	// Register messageCreate as a callback for the messageCreate events.
	discord.AddHandler(messageCreate)

	// Register ready as a callback for the ready events.
	discord.AddHandler(ready)

	// Open the websocket and begin listening.
	err = discord.Open()
	if err != nil {
		log.Fatal("Error opening Discord session: ", err)
	}

	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})

}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	log.Info("HPBot is now running. Press CTRL-C to exit.")
}

type GuildRoles map[string]string

var rolesMap = make(map[string]GuildRoles)

func getRoleName(s *discordgo.Session, roleId string, guildID string) (string, error) {
	if rolesMap[guildID] != nil {
		return rolesMap[guildID][roleId], nil
	}

	fmt.Println(guildID)

	guildRoles, err := s.GuildRoles(guildID)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	rolesMap[guildID] = make(GuildRoles)
	for _, role := range guildRoles {
		rolesMap[guildID][role.ID] = role.Name
	}

	return rolesMap[guildID][roleId], nil
}

func getHouseForMember(s *discordgo.Session, member *discordgo.Member, guildID string) (string, error) {
	for _, role := range member.Roles {
		log.WithFields(log.Fields{
			"guildID": guildID,
			"role":    role,
		}).Info("Getting house for member with")

		role, err := getRoleName(s, role, guildID)
		if err != nil {
			return "", err
		}

		if _, ok := housePoints[role]; ok {
			return role, nil
		}
	}
	return "", errors.New("Role not found")
}

func canWeGivePoints(s *discordgo.Session,
	giver *discordgo.Member,
	receiver *discordgo.Member,
	receiverHouse string,
	guildID string,
) bool {

	if giver.User.ID == receiver.User.ID {
		// return false
	}

	giverHouse := getHouseForMember(s, giver, guildID)
	log.Debug("%s is a %s\n", giver.Nick, giverHouse)

	return true
	// return giverHouse != receiverHouse
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("--------------------------------------------------")
	fmt.Println(m.Content)
	if len(m.Mentions) == 0 {
		fmt.Println("No mentions, return early")
		return
	}
	if !strings.Contains(m.Content, "++") {
		fmt.Println("Didn't plus plus, return early")
		return
	}

	fmt.Println(m.Author)
	channel, _ := s.Channel(m.ChannelID)
	guildID := channel.GuildID
	giver, _ := s.GuildMember(guildID, m.Author.ID)
	receiver, _ := s.GuildMember(guildID, m.Mentions[0].ID)

	receiverHouse := getHouseForMember(s, receiver, guildID)
	fmt.Printf("%s is a %s\n", receiver.Nick, receiverHouse)

	if !canWeGivePoints(s, giver, receiver, receiverHouse, guildID) {
		fmt.Println("cannot give points")
		return
	}
	addPoints(m.Mentions[0].ID, guildID, receiverHouse)
}
