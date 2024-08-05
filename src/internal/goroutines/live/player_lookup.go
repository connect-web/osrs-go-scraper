package live

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/goroutines"
	"github.com/connect-web/Low-Latency/internal/metrics"
	"github.com/connect-web/Low-Latency/internal/requests/players"
	"github.com/connect-web/Low-Latency/internal/utility/entities"
	"github.com/connect-web/Low-Latency/internal/utility/playerutils"
	"time"
)

func PlayersLookup(playerMapChunk map[string]int, proxyIterator *entities.ProxyIterator) goroutines.PlayerLookupResults {
	results := goroutines.PlayerLookupResults{
		Players:      make([]playerutils.SimplePlayer, 0),
		PlayerTotals: make([]playerutils.PlayerTotals, 0),
		NotFound:     make(map[int]struct{}),
	}

	for username, playerId := range playerMapChunk {
		placeholder := playerutils.NewSimplePlayer(playerId, username)
		player, err := players.GetPlayerStats(placeholder, proxyIterator)
		if err == nil {
			metrics.StatsFound++

			player.Calculations()

			playerTotals := playerutils.NewPlayerTotals()
			playerTotals.Calculate(player)

			currentTimestamp := time.Now()

			// Set last updated & playerId to the same.
			player.LastUpdated = currentTimestamp

			playerTotals.PlayerId = player.PID
			playerTotals.LastUpdated = currentTimestamp

			if player.PID == 0 {
				fmt.Println(player.PID)
			}

			results.Players = append(results.Players, player)
			results.PlayerTotals = append(results.PlayerTotals, playerTotals)

		} else if err.Error() == "page not found" {
			// user not found
			results.NotFound[playerId] = struct{}{}
			continue
		}
	}
	return results
}
