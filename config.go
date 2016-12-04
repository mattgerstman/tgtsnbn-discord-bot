package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	DiscordToken string          `json:"discord_token"`
	Houses       map[string]bool `json:"houses"`
}

var config *Configuration

/**
 * Loads config variables from file into Config struct.
 */
func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Failed to load config with message: ", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("Failed to decode config with message: ", err)
	}

	// Read this from the environment.
	config.DiscordToken = os.Getenv("DISCORD_TOKEN")
}

/**
 * Gets config struct. Initializes it if necessary.
 */
func GetConfig() *Configuration {
	if config == nil {
		loadConfig()
	}
	return config
}
