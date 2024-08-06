package main

import (
	"flag"
	"fmt"
	gainDb "github.com/connect-web/Low-Latency/internal/db/gain"
	"github.com/connect-web/Low-Latency/internal/goroutines/gain"

	"time"
)

var (
	threads                   = 75
	insertPlayerSize          = 10_000
	insertNotFoundSize        = 5_000
	queryFetchSize            = 100_000
	maxNbConcurrentGoroutines = flag.Int("MaxRoutines", threads, "The number of goroutines that are allowed to run concurrently")
)

func main() {
	for {
		players, NewPlayerError := gainDb.GetPlayersRequiringGains(queryFetchSize)
		fmt.Printf("%d players require gains\n", len(players))
		if len(players) == 0 {
			time.Sleep(30 * time.Minute)
		}
		if NewPlayerError == nil {
			gain.FindPlayers(players, maxNbConcurrentGoroutines, insertPlayerSize, insertNotFoundSize)
		}
	}
}
