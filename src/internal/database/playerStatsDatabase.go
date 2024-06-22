package database

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utils/stats"
)

/*
Contains functions with short-lived database connections
defer closes the connection after function has finished
*/

func GetNewPlayers(limit int) (map[string]int, error) {
	dbClient := NewDBClient()
	err := dbClient.Connect()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return map[string]int{}, err
	}
	defer func(dbClient *DBClient) {
		dbClientErr := dbClient.Close()
		if dbClientErr != nil {
			fmt.Println(dbClientErr.Error())
		}
	}(dbClient)

	Client := NewLeaderboardClient(dbClient)

	playerNameToId, err := Client.fetchUnknownStats(limit)

	if err == nil {
		fmt.Printf("Successfully found %d new users to scrape.\n", len(playerNameToId))
	}
	return playerNameToId, err
}

func GetOutdatedPlayers(limit int) (map[string]int, error) {
	dbClient := NewDBClient()
	err := dbClient.Connect()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return map[string]int{}, err
	}
	defer func(dbClient *DBClient) {
		dbClientErr := dbClient.Close()
		if dbClientErr != nil {
			fmt.Println(dbClientErr.Error())
		}
	}(dbClient)

	Client := NewLeaderboardClient(dbClient)

	playerNameToId, err := Client.fetchOutdatedStats(limit)

	if err == nil {
		fmt.Printf("Successfully found %d new users to scrape.\n", len(playerNameToId))
	}
	return playerNameToId, err
}

func InsertSimplePlayers(players []stats.SimplePlayer) error {
	dbClient := NewDBClient()
	err := dbClient.Connect()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return err
	}
	defer func(dbClient *DBClient) {
		dbClientErr := dbClient.Close()
		if dbClientErr != nil {
			fmt.Println(dbClientErr.Error())
		}
	}(dbClient)

	Client := NewLeaderboardClient(dbClient)

	err = Client.insertOrUpdatePlayers(players)

	if err == nil {
		fmt.Printf("Successfully inserted or updated %d players.\n", len(players))
	}
	return err
}

func InsertNotFound(players map[int]struct{}) error {
	dbClient := NewDBClient()
	err := dbClient.Connect()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return err
	}
	defer func(dbClient *DBClient) {
		dbClientErr := dbClient.Close()
		if dbClientErr != nil {
			fmt.Println(dbClientErr.Error())
		}
	}(dbClient)

	Client := NewLeaderboardClient(dbClient)
	err = Client.insertNotFound(players)

	if err == nil {
		fmt.Printf("Successfully inserted %d players to not found.\n", len(players))
	}
	return err
}
