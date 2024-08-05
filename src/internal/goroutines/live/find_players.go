package live

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/db/live"
	"github.com/connect-web/Low-Latency/internal/goroutines"
	"github.com/connect-web/Low-Latency/internal/metrics"
	"github.com/connect-web/Low-Latency/internal/utility/entities"
	"github.com/connect-web/Low-Latency/internal/utility/playerutils"
	"sync"
)

func FindPlayerLive(playerMap map[string]int, maxNbConcurrentGoroutines *int, insertSizePlayers int, insertSizeNotFound int) {

	var players []playerutils.SimplePlayer
	var playerTotals []playerutils.PlayerTotals
	notFoundPlayerIds := make(map[int]struct{}) // same as array but without duplicates.

	var wg sync.WaitGroup
	concurrentGoroutines := make(chan struct{}, *maxNbConcurrentGoroutines)
	respChan := make(chan goroutines.PlayerLookupResults)

	chunkedPlayerMap := entities.ChunkUserMap(playerMap, 10)

	for _, playerChunk := range chunkedPlayerMap {
		wg.Add(1)
		go func(playerMapChunk map[string]int, c chan goroutines.PlayerLookupResults, iterator *entities.ProxyIterator) {
			defer wg.Done()
			concurrentGoroutines <- struct{}{}
			c <- PlayersLookup(playerMapChunk, iterator)
			<-concurrentGoroutines
		}(playerChunk, respChan, entities.ProxyList)
	}

	for range chunkedPlayerMap {
		r := <-respChan

		for _, player := range r.Players {
			players = append(players, player)
			if len(players)%insertSizePlayers == 0 {
				players = insertPlayers(players)
			}
		}

		for _, playerTotal := range r.PlayerTotals {
			playerTotals = append(playerTotals, playerTotal)
			if len(playerTotals)%insertSizePlayers == 0 {
				playerTotals = insertPlayerTotals(playerTotals)
			}
		}

		for playerId := range r.NotFound {
			notFoundPlayerIds[playerId] = struct{}{}
			if len(notFoundPlayerIds)%insertSizeNotFound == 0 {
				notFoundPlayerIds = insertNotFound(notFoundPlayerIds)
			}
		}
	}
	wg.Wait()
	if len(notFoundPlayerIds) > 0 {
		insertNotFound(notFoundPlayerIds)
	}
	if 0 < len(players) {
		insertPlayers(players)
	}
}

func insertNotFound(notFoundPlayerIds map[int]struct{}) map[int]struct{} {
	err := live.InsertNotFound(notFoundPlayerIds)

	if err != nil {
		fmt.Printf("Failed to submit %d stats.\n", len(notFoundPlayerIds))
		return notFoundPlayerIds
	}
	metrics.NotfoundInserted += len(notFoundPlayerIds)
	fmt.Printf("%.2fK Not found @ %.2fK/hr\n", float64(metrics.NotfoundInserted)/1000, metrics.GetHourly(metrics.NotfoundInserted))
	return map[int]struct{}{} // Reset to prevent duplicate entries.
}

func insertPlayers(players []playerutils.SimplePlayer) []playerutils.SimplePlayer {
	err := live.InsertSimplePlayers(players)
	if err != nil {
		fmt.Printf("Failed to submit %d player stats.\n", len(players))
		return players
	}
	metrics.PlayersInserted += len(players)
	fmt.Printf("%.2fK Players stats inserted @ %.2fK/hr\n", float64(metrics.PlayersInserted)/1000, metrics.GetHourly(metrics.PlayersInserted))

	return []playerutils.SimplePlayer{} // Reset to prevent duplicate entries.
}

func insertPlayerTotals(players []playerutils.PlayerTotals) []playerutils.PlayerTotals {
	err := live.InsertPlayerLiveStats(players)
	if err != nil {
		fmt.Printf("Failed to submit %d player combat stats.\n", len(players))
		return players
	}
	metrics.PlayerTotals += len(players)
	fmt.Printf("%.2fK Player total's stats inserted @ %.2fK/hr\n", float64(metrics.PlayerTotals)/1000, metrics.GetHourly(metrics.PlayerTotals))

	return []playerutils.PlayerTotals{} // Reset to prevent duplicate entries.
}
