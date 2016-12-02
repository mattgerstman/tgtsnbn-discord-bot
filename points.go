package main

import log "github.com/Sirupsen/logrus"

type Users struct {
}

type HousePoint struct {
	House     string `json:"house"`
	NumPoints string `json:"num_points"`
}

func canWeGivePoints(
	giverID string,
	giverHouse string,
	receiverID string,
	receiverHouse string,
) bool {

	if giverID == receiverID {
		// return false
	}

	return true
	// return giverHouse != receiverHouse
}

func addPoints(userID string, guildID string, house string) *ApplicationError {
	houseMap := GetHouseMap()
	if _, ok := houseMap[house]; !ok {
		log.WithFields(log.Fields{
			"house": house,
		}).Warn("Invalid house")
		return NewApplicationErrorWithoutError("Invalid House", ErrorInvalidHouse)
	}

	_, err := db.Exec(
		`INSERT INTO users (user_id, guild_id, house, num_points) VALUES (?,?,?,?)
		ON DUPLICATE KEY UPDATE num_points = num_points + 10, house = ?`,
		userID, guildID, house, 10, house,
	)
	if err != nil {
		return NewApplicationError("Error adding points", err, ErrorDatabase)
	}

	getHouseStandings(guildID)
	return nil
}

func getPointsForUser(userID string, guildID string) {
	db.QueryRow(`SELECT num_points FROM users WHERE user_id = ? AND guild_id = ?`, userID, guildID)
}

func getHouseStandings(guildID string) []HousePoint {
	housePoints := make([]HousePoint, 0)

	rows, err := db.Query(`SELECT house, SUM(num_points) as num_points FROM users WHERE guild_id = ? GROUP BY house ORDER BY num_points DESC`, guildID)
	if err != nil {
		log.Error(err)
		return nil
	}

	for rows.Next() {
		var housePoint HousePoint
		rows.Scan(&housePoint.House, &housePoint.NumPoints)
		log.WithFields(log.Fields{
			"house":      housePoint.House,
			"num_points": housePoint.NumPoints,
		}).Info("Points for house")
		housePoints = append(housePoints, housePoint)
	}

	return housePoints
}
