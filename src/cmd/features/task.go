package main

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/database"
	"github.com/connect-web/Low-Latency/internal/utils/stats"
	"log"
	"time"
)

func processPlayers(players []stats.SimplePlayer) []stats.AdvancedPlayer {
	advancedPlayers := []stats.AdvancedPlayer{}
	for _, player := range players {
		var advancedPlayer = stats.AdvancedPlayer{}
		advancedPlayer.Calculate(player)
		advancedPlayers = append(advancedPlayers, advancedPlayer)
	}
	return advancedPlayers
}

func main() {
	dbClient := database.NewDBClient()
	err := dbClient.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database client. %s\n", err.Error())
	}
	Client := database.NewLeaderboardClient(dbClient)

	for {
		finished := processFeatures(Client.FetchUnknownFeatures, Client.InsertOrUpdateFeatures, 100_000)
		if finished {
			updatesFinished := processFeatures(Client.FetchOutdatedFeatures, Client.InsertOrUpdateFeatures, 100_000)
			if updatesFinished {
				time.Sleep(10 * time.Minute)
			}
		}
	}

}

func processFeatures(
	playerFetcher func(int) ([]stats.SimplePlayer, error),
	insert func([]stats.AdvancedPlayer) error,
	limit int) (noPlayers bool) {

	players, err := playerFetcher(limit)
	if err != nil {
		log.Printf("Failed to get players requiring features: %s\n", err.Error())
		return false
	}

	if len(players) == 0 {
		return true
	}

	now := time.Now().UnixMilli()

	advancedPlayerList := processPlayers(players)

	msDiff := time.Now().UnixMilli() - now
	// 0.12 seconds for 100k players. Goroutines are not required here...
	fmt.Printf("Processed %d players in %.4fs\n", len(players), float64(msDiff)/1000)

	insertErr := insert(advancedPlayerList)
	if insertErr != nil {
		log.Printf("Failed to insert advanced players %s\n", insertErr.Error())
	}
	return false
}
