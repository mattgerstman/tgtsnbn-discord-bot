package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/mattgerstman/discordgo"
)

func main() {

	// Initialize the database.
	log.Info("Connecting to database...")
	initDB()
	log.Info("Connected to database")

	// Load the config.
	log.Info("Loading config file...")
	config := GetConfig()
	log.Info("Loaded config file")

	log.Info("Connecting to Discord...")
	// Log into discord.
	discord, err := discordgo.New(config.Username, config.Password)
	if err != nil {
		log.Fatal("Failed to connect to discord ", err)
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
	log.Info("TGTSNBN Discord Bot is now running. Press CTRL-C to exit.")
}
