package gain

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/db/gain"
	"github.com/connect-web/Low-Latency/internal/db/live"
	"github.com/connect-web/Low-Latency/internal/goroutines"
	"github.com/connect-web/Low-Latency/internal/metrics"
	"github.com/connect-web/Low-Latency/internal/utility/entities"
	"github.com/connect-web/Low-Latency/internal/utility/playerutils"
	"sync"
	"time"
)

func FindPlayers(players []playerutils.SimplePlayer, maxNbConcurrentGoroutines *int, insertSizePlayers int, insertSizeNotFound int) {
	var playerDifferences []playerutils.SimplePlayer
	notFoundPlayerIds := make(map[int]struct{}) // same as array but without duplicates.

	chunkedPlayerMap := entities.ChunkSimplePlayer(players, 10)

	concurrentGoroutines := make(chan struct{}, *maxNbConcurrentGoroutines)
	respChan := make(chan goroutines.PlayerLookupResults)
	var wg sync.WaitGroup

	for _, playerChunk := range chunkedPlayerMap {
		wg.Add(1)
		go func(playerMapChunk []playerutils.SimplePlayer, c chan goroutines.PlayerLookupResults, iterator *entities.ProxyIterator) {
			defer wg.Done()
			concurrentGoroutines <- struct{}{}
			c <- PlayersLookup(playerMapChunk, iterator)
			<-concurrentGoroutines
		}(playerChunk, respChan, entities.ProxyList)
	}

	for range chunkedPlayerMap {
		r := <-respChan
		for _, player := range r.Players {
			playerDifferences = append(playerDifferences, player)

			if len(playerDifferences)%insertSizePlayers == 0 {
				playerDifferences = publishGains(playerDifferences)
			}
		}

		for playerId := range r.NotFound {
			notFoundPlayerIds[playerId] = struct{}{}

			if len(notFoundPlayerIds)%insertSizeNotFound == 0 {
				notFoundPlayerIds = publishNotFound(notFoundPlayerIds)
			}
		}
	}
	// wait for all routines to finish.
	wg.Wait()

	// final insert if data to insert.
	if len(notFoundPlayerIds) > 0 {
		notFoundPlayerIds = publishNotFound(notFoundPlayerIds)
	}

	for 0 < len(playerDifferences) {
		playerDifferences = publishGains(playerDifferences)
		time.Sleep(10 * time.Second)
	}
}

func publishNotFound(notFoundPlayerIds map[int]struct{}) map[int]struct{} {
	err := live.InsertNotFound(notFoundPlayerIds)
	if err != nil {
		fmt.Printf("Failed to submit %d stats.\n", len(notFoundPlayerIds))
		return notFoundPlayerIds
	}

	metrics.NotfoundInserted += len(notFoundPlayerIds)
	fmt.Printf("%.2fK Not found @ %.2fK/hr\n", float64(metrics.NotfoundInserted)/1000, metrics.GetHourly(metrics.NotfoundInserted))

	return map[int]struct{}{} // Reset to prevent duplicate entries.
}

func publishGains(playerDifferences []playerutils.SimplePlayer) []playerutils.SimplePlayer {
	err := gain.PublishGains(playerDifferences)
	if err != nil {
		fmt.Printf("Failed to submit %d stats.\n", len(playerDifferences))
		return playerDifferences
	}
	metrics.PlayersInserted += len(playerDifferences)
	fmt.Printf("%.2fK Players stats inserted @ %.2fK/hr\n", float64(metrics.PlayersInserted)/1000, metrics.GetHourly(metrics.PlayersInserted))

	return []playerutils.SimplePlayer{} // Reset to prevent duplicate entries.
}
