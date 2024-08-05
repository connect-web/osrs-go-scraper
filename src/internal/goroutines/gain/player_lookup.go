package gain

import (
	"github.com/connect-web/Low-Latency/internal/goroutines"
	"github.com/connect-web/Low-Latency/internal/metrics"
	"github.com/connect-web/Low-Latency/internal/requests/players"
	"github.com/connect-web/Low-Latency/internal/utility/entities"
	"github.com/connect-web/Low-Latency/internal/utility/playerutils"
	"log"
)

func PlayersLookup(playerChunk []playerutils.SimplePlayer, proxyIterator *entities.ProxyIterator) goroutines.PlayerLookupResults {
	results := goroutines.PlayerLookupResults{
		Players:  []playerutils.SimplePlayer{},
		NotFound: make(map[int]struct{}),
	}

	for _, oldPlayer := range playerChunk {
		placeholderPlayer := playerutils.SimplePlayer{
			Username:    oldPlayer.Username,
			LastUpdated: oldPlayer.LastUpdated,
			PID:         oldPlayer.PID,
		}
		player, err := players.GetPlayerStats(placeholderPlayer, proxyIterator)

		if err == nil {
			if len(player.Skills) == 0 && len(player.Minigames) == 0 {
				continue // skip if the skills & minigames was empty from API
			}
			metrics.StatsFound++
			playerGains := oldPlayer.CalculateGains(player)
			results.Players = append(results.Players, playerGains)

		} else if err.Error() == "page not found" {
			// user not found
			results.NotFound[oldPlayer.PID] = struct{}{}
			continue
		} else {
			log.Printf("PlayersLookup: %s\n", err.Error())
		}
	}
	return results
}
