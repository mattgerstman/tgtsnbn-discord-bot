package main

import log "github.com/Sirupsen/logrus"

type Users struct {
}

/**
 * Adds points to a user/house, persists to database.
 */
func AddPoints(userID string, guildID string, house string) *ApplicationError {
	houseMap := GetHouseMap()
	if _, ok := houseMap[house]; !ok {
		log.WithFields(log.Fields{
			"house": house,
		}).Warn("Invalid house")
		return NewApplicationErrorWithoutError("Invalid House", ErrorInvalidHouse)
	}

	db := GetDB()
	_, err := db.Exec(
		`INSERT INTO users (user_id, guild_id, house, num_points)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, guild_id)
		DO UPDATE SET num_points = users.num_points + $4, house = $3
		WHERE users.user_id = $1 AND users.guild_id = $2`,
		userID, guildID, house, 10,
	)
	if err != nil {
		return NewApplicationError("Error adding points", err, ErrorDatabase)
	}

	return nil
}

/**
 * Returns the number of points a user has.
 */
func GetPointsForUser(
	userID string,
	guildID string,
) (numPoints int, appErr *ApplicationError) {
	db := GetDB()
	err := db.QueryRow(
		`SELECT num_points FROM users WHERE user_id = $1 AND guild_id = $2`,
		userID,
		guildID,
	).Scan(&numPoints)

	if err != nil {
		return 0,
			NewApplicationError("Error fetching points for user", err, ErrorDatabase)
	}
	return numPoints, nil
}

type HousePoint struct {
	House     string `json:"house"`
	NumPoints string `json:"num_points"`
}

/**
 * Returns a list of all house standings.
 */
func GetHouseLeaderboard(guildID string) ([]HousePoint, *ApplicationError) {
	housePoints := make([]HousePoint, 0)

	db := GetDB()
	rows, err := db.Query(
		`SELECT house, SUM(num_points) as num_points FROM users WHERE guild_id = $1
		GROUP BY house ORDER BY num_points DESC`,
		guildID,
	)
	if err != nil {
		return nil, NewApplicationError(
			"Unable to fetch house points from db",
			err,
			ErrorDatabase,
		)
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

	return housePoints, nil
}

/**
 * Returns the number of points a single house has.
 */
func GetPointsForHouse(house string, guildID string) (numPoints int, appErr *ApplicationError) {
	db := GetDB()
	err := db.QueryRow(
		`SELECT SUM(num_points) FROM users WHERE house = $1 AND guild_id = $2`,
		house,
		guildID,
	).Scan(&numPoints)

	if err != nil {
		return 0,
			NewApplicationError("Error fetching points for house", err, ErrorDatabase)
	}
	return numPoints, nil
}
