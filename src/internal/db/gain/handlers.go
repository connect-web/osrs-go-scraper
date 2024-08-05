package gain

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utility/playerutils"
)

func GetPlayersRequiringGains(limit int) ([]playerutils.SimplePlayer, error) {
	Client := NewGainsClient()
	defer Client.Close()

	players, err := Client.fetchPlayersRequireGains(limit)

	if err == nil {
		fmt.Printf("Successfully found %d new users to scrape gains.\n", len(players))
	}
	return players, err
}

func PublishGains(simplePlayers []playerutils.SimplePlayer) error {
	Client := NewGainsClient()
	defer Client.Close()

	err := Client.InsertOrUpdateGains(simplePlayers)

	if err == nil {
		fmt.Printf("Successfully inserted %d players to not found.\n", len(simplePlayers))
	}
	return err
}
