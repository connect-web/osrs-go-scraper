package name

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/db/name"
	"github.com/connect-web/Low-Latency/internal/goroutines"
	"github.com/connect-web/Low-Latency/internal/metrics"
	"github.com/connect-web/Low-Latency/internal/requests/leaderboards"
	"github.com/connect-web/Low-Latency/internal/utility/entities"
	"github.com/connect-web/Low-Latency/internal/utility/nameutils"
	"sync"
	"time"
)

func Run(maxNbConcurrentGoroutines *int, usernameBatchSize int, low_memory bool) {
	limitManager := FindLimits(maxNbConcurrentGoroutines)
	/*
		Depreciated method
		limitManager := name.NewPageLimitManager()
	*/

	concurrentGoroutines := make(chan struct{}, *maxNbConcurrentGoroutines)
	respChan := make(chan goroutines.AccInfo)
	var wg sync.WaitGroup

	var newNames = make(map[string]struct{}) // same as array but without duplicates.

	for _, hiscoreTableInfo := range limitManager.Limits {

		for page := 0; page < hiscoreTableInfo.Limit; page++ {
			wg.Add(1)
			go func(hiscoreType nameutils.HiscoreType, page int, c chan goroutines.AccInfo, iterator *entities.ProxyIterator) {
				defer wg.Done()
				concurrentGoroutines <- struct{}{}
				uniqueNames, _ := leaderboards.GetNames(hiscoreType, page, iterator)
				c <- goroutines.AccInfo{Unique_names: uniqueNames}
				<-concurrentGoroutines
			}(hiscoreTableInfo, page, respChan, entities.ProxyList)
		}
	}

	for _, hiscoreTableInfo := range limitManager.Limits {
		for page := 0; page < hiscoreTableInfo.Limit; page++ {
			r := <-respChan
			for username := range r.Unique_names {
				if username == "" {
					continue
				}
				if !low_memory && name.IsUsernameKnown(username) {
					continue // if High memory and username is already in memory then continue
				}
				newNames[username] = struct{}{}
				if len(newNames)%usernameBatchSize == 0 {
					newNames = SubmitUsernames(newNames)
				}
			}
		}
	}
	// wait for all routines to finish.
	wg.Wait()

	// final insert if data to insert.
	for 0 < len(newNames) {
		newNames = SubmitUsernames(newNames)
		time.Sleep(10 * time.Second)
	}
}

func SubmitUsernames(newNames map[string]struct{}) map[string]struct{} {
	err := name.SubmitUsernames(newNames)
	if err != nil {
		fmt.Printf("Failed to submit %d usernames.\n", len(newNames))
		return newNames
	}
	metrics.UsernamesInserted += len(newNames)
	fmt.Printf("%.2fK Usernames inserted in database!\n", float64(metrics.UsernamesInserted/1000))

	fmt.Printf("%.2fK unique Names @ %.2fK/hr\n", float64(metrics.UsernamesInserted)/1000, metrics.GetHourly(metrics.UsernamesInserted))

	return map[string]struct{}{} // Reset to prevent duplicate entries.
}
