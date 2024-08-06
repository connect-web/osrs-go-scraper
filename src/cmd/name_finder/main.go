package main

import (
	"flag"
	nameDb "github.com/connect-web/Low-Latency/internal/db/name"
	"github.com/connect-web/Low-Latency/internal/goroutines/name"
)

var (
	threads                   = 75
	usernameBatchSize         = 5_000
	LowMemory                 = false
	maxNbConcurrentGoroutines = flag.Int("MaxRoutines", threads, "The number of goroutines that are allowed to run concurrently")
)

func main() {
	if !LowMemory {
		nameDb.LoadKnownUsernames()
	}
	for {
		name.Run(maxNbConcurrentGoroutines, usernameBatchSize, LowMemory)
		// todo add a log to see how long it takes to scrape the full Hiscores.
		// new table in database for logging
	}
}
