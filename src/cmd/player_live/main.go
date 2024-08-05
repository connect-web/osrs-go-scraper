package main

import (
	"flag"
	liveDb "github.com/connect-web/Low-Latency/internal/db/live"
	"github.com/connect-web/Low-Latency/internal/goroutines/live"
	"time"
)

var (
	threads = 250

	queryFetchSize            = 100_000
	insertSizePlayers         = 10_000
	insertSizeNotFound        = 10_000
	maxNbConcurrentGoroutines = flag.Int("MaxRoutines", threads, "The number of goroutines that are allowed to run concurrently")
)

func main() {

	for {
		playerMap, NewPlayerError := liveDb.GetNewPlayers(queryFetchSize)
		if len(playerMap) == 0 {
			time.Sleep(10 * time.Minute)
		}
		if NewPlayerError == nil {
			live.FindPlayerLive(playerMap, maxNbConcurrentGoroutines, insertSizePlayers, insertSizeNotFound)
		}
	}
}
