package main

import "github.com/mattgerstman/discordgo"

/*
 * Helper to get the name for a member.
 */
func GetNameForMember(member *discordgo.Member) string {
	if member.Nick != "" {
		return member.Nick
	}

	return member.User.Username
}
