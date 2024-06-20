package main

import (
	"flag"
	"fmt"
	database "github.com/connect-web/Low-Latency-DB"
	utils "github.com/connect-web/Low-Latency-Utils"
	"log"
	"sync"
	"time"
)

var (
	start                     = time.Now().Unix()
	threads                   = 50
	maxNbConcurrentGoroutines = flag.Int("MaxRoutines", threads, "The number of goroutines that are allowed to run concurrently")
	proxyIterator             = utils.NewProxyIterator("proxies.txt")
	statsFound                = 0
	playersInserted           = 0
	notfoundInserted          = 0

	insertSizeNotFound = 500 // larger transaction will have less overhead
	insertSizePlayers  = 3000
)

type PlayerLookupResults struct {
	Players  []utils.SimplePlayer
	NotFound map[int]struct{}
}

func main() {
	for {
		// if finished getting all new players then update old ones
		if findNewPlayers() {
			updatePlayers()
		}
	}
}

func findNewPlayers() bool { // returns true if finished getting all new players
	playerMap, NewPlayerError := database.GetNewPlayers(50_000)

	if len(playerMap) == 0 {
		return true
	}
	if NewPlayerError == nil {
		FindPlayers(playerMap)
	}
	return false
}

func updatePlayers() {
	playerMap, NewPlayerError := database.GetOutdatedPlayers(5000)
	if NewPlayerError == nil {
		FindPlayers(playerMap)
	}
}

func FindPlayers(playerMap map[string]int) {
	players := []utils.SimplePlayer{}
	notFoundPlayerIds := make(map[int]struct{}) // same as array but without duplicates.

	chunkedPlayerMap := utils.ChunkUserMap(playerMap, 10)

	concurrentGoroutines := make(chan struct{}, *maxNbConcurrentGoroutines)
	respChan := make(chan PlayerLookupResults)
	var wg sync.WaitGroup

	for _, playerChunk := range chunkedPlayerMap {
		wg.Add(1)
		go func(playerMapChunk map[string]int, c chan PlayerLookupResults, iterator *utils.ProxyIterator) {
			defer wg.Done()
			concurrentGoroutines <- struct{}{}
			c <- PlayersLookup(playerMapChunk, iterator)
			<-concurrentGoroutines
		}(playerChunk, respChan, proxyIterator)
	}

	for range chunkedPlayerMap {
		r := <-respChan

		for _, player := range r.Players {
			players = append(players, player)

			if len(players)%insertSizePlayers == 0 {
				err := database.InsertSimplePlayers(players)
				if err != nil {
					fmt.Printf("Failed to submit %d usernames.\n", len(players))
				} else {
					playersInserted += len(players)
					fmt.Printf("%.2fK unique Names @ %.2fK/hr\n", float64(playersInserted)/1000, getHourly(playersInserted))
					// Reset to prevent duplicate entries.
					players = []utils.SimplePlayer{}
				}
			}
		}

		for playerId := range r.NotFound {
			notFoundPlayerIds[playerId] = struct{}{}

			if len(notFoundPlayerIds)%insertSizeNotFound == 0 {
				err := database.InsertNotFound(notFoundPlayerIds)
				if err != nil {
					fmt.Printf("Failed to submit %d usernames.\n", len(notFoundPlayerIds))
				} else {
					notfoundInserted += len(notFoundPlayerIds)
					fmt.Printf("%.2fK Not found @ %.2fK/hr\n", float64(notfoundInserted)/1000, getHourly(notfoundInserted))
					// Reset to prevent duplicate entries.
					notFoundPlayerIds = map[int]struct{}{}
				}
			}
		}

	}
	// wait for all routines to finish.
	wg.Wait()
	// final insert if data to insert.

	if len(notFoundPlayerIds) > 0 {
		err := database.InsertNotFound(notFoundPlayerIds)
		if err != nil {
			fmt.Printf("Failed to submit %d usernames.\n", len(notFoundPlayerIds))
		} else {
			notfoundInserted += len(notFoundPlayerIds)
			fmt.Printf("%.2fK Not found @ %.2fK/hr\n", float64(notfoundInserted)/1000, getHourly(notfoundInserted))
			// Reset to prevent duplicate entries.
			notFoundPlayerIds = map[int]struct{}{}
		}
	}

	if len(players) > 0 {
		err := database.InsertSimplePlayers(players)
		if err != nil {
			fmt.Printf("Failed to submit %d usernames.\n", len(players))
		} else {
			playersInserted += len(players)
			fmt.Printf("%.2fK unique Names @ %.2fK/hr\n", float64(playersInserted)/1000, getHourly(playersInserted))
			// Reset to prevent duplicate entries.
			players = []utils.SimplePlayer{}
		}
	}

}

func PlayersLookup(playerMapChunk map[string]int, proxyIterator *utils.ProxyIterator) PlayerLookupResults {
	results := PlayerLookupResults{
		Players:  []utils.SimplePlayer{},
		NotFound: make(map[int]struct{}),
	}

	for username, player_id := range playerMapChunk {
		player, err := utils.Get_player_stats(username, player_id, proxyIterator)

		if err == nil {
			statsFound++
			// shows player stats
			//fmt.Printf("%d: %v\n", player)
			// fmt.Println(statsFound)
			results.Players = append(results.Players, player)

		} else if err.Error() == "Page not found" {
			// user not found
			results.NotFound[player_id] = struct{}{}
			continue
		} else {
			log.Printf(err.Error())
		}
	}
	return results
}

func test() {
	playerMap, NewPlayerError := database.GetNewPlayers(20)

	if NewPlayerError != nil {
		log.Fatalf("Failed to get new players: %s\n", NewPlayerError.Error())
	}

	chunkedPlayerMap := utils.ChunkUserMap(playerMap, 10)

	for _, playerMapChunk := range chunkedPlayerMap {
		PlayersLookup(playerMapChunk, proxyIterator)
	}
}

func getHourly(size int) float64 {
	secondsRan := time.Now().Unix() - start
	PerSecond := float64(size) / float64(secondsRan)
	PerHour := (PerSecond * 3600) / 1000 // K players per hour.
	return PerHour
}
