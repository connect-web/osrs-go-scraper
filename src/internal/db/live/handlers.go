package live

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utility/playerutils"
)

/*
Contains functions with short-lived database connections
defer closes the connection after function has finished
*/

func GetNewPlayers(limit int) (map[string]int, error) {
	Client := NewLiveClient()
	defer Client.Close()
	playerNameToId, err := Client.fetchUnknownStats(limit)

	if err == nil {
		fmt.Printf("Successfully found %d new users to scrape.\n", len(playerNameToId))
	}
	return playerNameToId, err
}

func InsertSimplePlayers(players []playerutils.SimplePlayer) error {
	Client := NewLiveClient()
	defer Client.Close()

	err := Client.insertOrUpdatePlayers(players)

	if err == nil {
		fmt.Printf("Successfully inserted or updated %d players.\n", len(players))
	}
	return err
}

func InsertNotFound(players map[int]struct{}) error {
	Client := NewLiveClient()
	defer Client.Close()

	err := Client.insertNotFound(players)

	if err == nil {
		fmt.Printf("Successfully inserted %d players to not found.\n", len(players))
	}
	return err
}

func InsertPlayerLiveStats(players []playerutils.PlayerTotals) error {
	Client := NewLiveClient()
	defer Client.Close()

	err := Client.InsertOrUpdateLiveStats(players)

	if err == nil {
		fmt.Printf("Successfully inserted or updated %d players.\n", len(players))
	}
	return err
}
