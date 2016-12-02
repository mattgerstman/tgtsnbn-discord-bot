package main

import log "github.com/Sirupsen/logrus"

type Users struct {
}

type HousePoint struct {
	House     string `json:"house"`
	NumPoints string `json:"num_points"`
}

func addPoints(userId string, guildID string, house string) error {
	houseMap := GetHouseMap()
	if _, ok := houseMap[house]; !ok {
		log.WithFields(log.Fields{
			"house": house,
		}).Warn("Invalid house")
		return nil
	}

	db := GetDB()
	_, err := db.Exec(
		`INSERT INTO users (user_id, guild_id, house, num_points) VALUES (?,?,?,?)
		ON DUPLICATE KEY UPDATE num_points = num_points + 10, house = ?`,
		userId, guildID, house, 10, house,
	)
	if err != nil {
		return err
	}

	getHouseStandings(guildID)
	return nil
}

func getPointsForUser(userId string, guildID string) {

}

func getHouseStandings(guildID string) []HousePoint {
	db := GetDB()

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
