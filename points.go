package main

import "fmt"

type Users struct {
}

type HousePoint struct {
	House     string `json:"house"`
	NumPoints string `json:"num_points"`
}

func addPoints(userId string, guildID string, house string) {
	houseMap := GetHouseMap()
	if _, ok := houseMap[house]; !ok {
		log.WithFields(log.Fields{
			"house": house,
		}).Warn("Invalid house")
		return
	}
	db := GetDB()

	_, err := db.Exec(`INSERT INTO users (user_id, guild_id, house, num_points) VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE num_points = num_points + 10, house = ?`, userId, guildID, house, 10, house)
	if err != nil {
		fmt.Println("ERROR executing")
		fmt.Println(err)
		return
	}

	getHouseStandings(guildID)
}

func getPointsForUser(userId string, guildID string) {

}

func getHouseStandings(guildID string) []HousePoint {
	db := GetDB()

	housePoints := make([]HousePoint, 0)

	rows, err := db.Query(`SELECT house, SUM(num_points) as num_points FROM users WHERE guild_id = ? GROUP BY house ORDER BY num_points DESC`, guildID)
	if err != nil {
		fmt.Println("Error fetching")
		fmt.Println(err)
		return nil
	}

	for rows.Next() {

		var housePoint HousePoint
		rows.Scan(&housePoint.House, &housePoint.NumPoints)
		fmt.Println(housePoint.House)
		fmt.Println(housePoint.NumPoints)
		fmt.Println(rows)
		housePoints = append(housePoints, housePoint)
	}

	fmt.Println(housePoints)

	return housePoints
}