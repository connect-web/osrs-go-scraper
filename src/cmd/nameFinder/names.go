package main

import (
	"flag"
	"fmt"
	"github.com/connect-web/Low-Latency/internal/database"
	"github.com/connect-web/Low-Latency/internal/utils/entities"
	"github.com/connect-web/Low-Latency/internal/utils/name"
	"sync"
	"time"
)

var (
	start                     = time.Now().Unix()
	threads                   = 50
	maxNbConcurrentGoroutines = flag.Int("MaxRoutines", threads, "The number of goroutines that are allowed to run concurrently")
	proxyIterator             = entities.NewProxyIterator("proxies.txt")
	usernames_verified        = 0
	usernameBatchSize         = 5000
)

type accInfo struct {
	unique_names map[string]bool
}

func main() {
	for {
		run()
		// todo add a log to see how long it takes to scrape the full Hiscores.
		// new table in database for logging
	}
}

func run() {
	//limitManager := name.NewPageLimitManager()
	limitManager := FindLimits()

	concurrentGoroutines := make(chan struct{}, *maxNbConcurrentGoroutines)
	respChan := make(chan accInfo)
	var wg sync.WaitGroup

	var new_names = make(map[string]struct{}) // same as array but without duplicates.

	for _, hiscore_table_info := range limitManager.Limits {

		for page := 0; page < hiscore_table_info.Limit; page++ {
			wg.Add(1)
			go func(hiscoreType name.HiscoreType, page int, c chan accInfo, iterator *entities.ProxyIterator) {
				/*
					Could add larger batch size to the goroutine however they are lightweight enough for current performance target
				*/
				defer wg.Done()
				concurrentGoroutines <- struct{}{}
				unique_names, _ := GetNames(hiscoreType, page, iterator)
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
				if len(new_names)%usernameBatchSize == 0 {
					// the Larger the batch size for usernames
					// = Less usernames due to removed duplicates
					err := database.SubmitUsernames(new_names)
					if err != nil {
						fmt.Printf("Failed to submit %d usernames.\n", len(new_names))
					} else {
						usernames_verified += len(new_names)
						fmt.Printf("%.2fK Usernames verified in database!\n", float64(usernames_verified/1000))

						secondsRan := time.Now().Unix() - start
						playersPerSecond := float64(usernames_verified) / float64(secondsRan)
						playersPerHour := (playersPerSecond * 3600) / 1000 // K players per hour.
						fmt.Printf("%.2fK unique Names @ %.2fK/hr\n", float64(usernames_verified)/1000, playersPerHour)

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
		err := database.SubmitUsernames(new_names)
		for err != nil {
			err = database.SubmitUsernames(new_names)
		}
		new_names = map[string]struct{}{}
	}
}
