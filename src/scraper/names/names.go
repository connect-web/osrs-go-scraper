package main

import (
	database "database/leaderboard.go"
	"flag"
	"fmt"
	"names.go/entities"
	"names.go/entities/limits"
	"names.go/requests"
	"sync"
	"time"
)

var (
	start                     = time.Now().Unix()
	threads                   = 50
	maxNbConcurrentGoroutines = flag.Int("MaxRoutines", threads, "The number of goroutines that are allowed to run concurrently")
	proxyIterator             = entities.NewProxyIterator("proxies.txt")
	usernames_verified        = 0
)

type accInfo struct {
	unique_names map[string]bool
}

func main() {
	run()
}

func run() {
	limitManager := FindLimits()

	concurrentGoroutines := make(chan struct{}, *maxNbConcurrentGoroutines)
	respChan := make(chan accInfo)
	var wg sync.WaitGroup

	var new_names = make(map[string]struct{}) // same as array but without duplicates.

	for _, hiscore_table_info := range limitManager.Limits {

		for page := 0; page < hiscore_table_info.Limit; page++ {
			wg.Add(1)
			go func(hiscoreType limits.HiscoreType, page int, c chan accInfo, iterator *entities.ProxyIterator) {
				/*
					Could add larger batch size to the goroutine however they are lightweight enough for current performance target
				*/
				defer wg.Done()
				concurrentGoroutines <- struct{}{}
				unique_names, _ := requests.GetNames(hiscoreType, page, iterator)
				c <- accInfo{unique_names: unique_names}
				<-concurrentGoroutines
			}(hiscore_table_info, page, respChan, proxyIterator)
		}
	}

	for _, hiscore_table_info := range limitManager.Limits {
		for page := 0; page < hiscore_table_info.Limit; page++ {
			r := <-respChan
			for username := range r.unique_names {
				if username == "" {
					continue
				}
				new_names[username] = struct{}{}

				if len(new_names)%1000 == 0 {
					secondsRan := time.Now().Unix() - start
					playersPerSecond := float64(len(new_names)) / float64(secondsRan)
					playersPerHour := (playersPerSecond * 3600) / 1000 // K players per hour.
					fmt.Printf("%.2fK unique Names @ %.2fK/hr\n", float64(len(new_names))/1000, playersPerHour)
				}
				if len(new_names)%100_000 == 0 {
					// the Larger the batch size for usernames
					// = Less usernames due to removed duplicates
					err := database.SubmitUsernames(new_names)
					if err != nil {
						fmt.Printf("Failed to submit %d usernames.\n", len(new_names))
					} else {
						usernames_verified += len(new_names)
						fmt.Printf("%.2fK Usernames verified in database!\n", float64(usernames_verified/1000))

						// Reset to prevent duplicate entries.
						new_names = map[string]struct{}{}
					}

				}
			}
		}
	}
	// wait for all routines to finish.
	wg.Wait()
	// final insert if data to insert.
	if len(new_names) > 0 {
		insertNewPlayers(new_names)
		new_names = map[string]struct{}{}
	}
}
