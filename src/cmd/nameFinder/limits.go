package main

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utils/entities"
	"github.com/connect-web/Low-Latency/internal/utils/name"
	"sync"
	"time"
)

/*
The old school runescape hiscores has a limit of the top 2 million users per hiscore ranking

Some hiscore rankings do not have 2 million users participating in the activity, therefore this finds the page ranges that the Leaderboard nameFinder will visit
- This prevents wasted requests on pages that do not exist.

*/

func FindLimits() *name.PageLimits {
	pageLimitManager := name.NewPageLimitManager()

	concurrentGoroutines := make(chan struct{}, *maxNbConcurrentGoroutines)
	respChan := make(chan name.PageLimitInfo)
	var wg sync.WaitGroup

	for _, hiscoreSkills := range name.HiscoreMinigames {
		wg.Add(1)

		go func(hiscoreType name.HiscoreType, page_limit int, c chan name.PageLimitInfo, iterator *entities.ProxyIterator) {
			defer wg.Done()
			concurrentGoroutines <- struct{}{}
			start_page := 10

			first_page, _ := GetNames(hiscoreType, 1, iterator)
			if len(first_page) == 0 {
				fmt.Printf("First page failed. %s\n", hiscoreType.Minigames)
			} else {
				delims := []int{5000, 1000, 100, 10}

				for _, delim_page := range delims {
					start_page = getValidPage(first_page, start_page, delim_page, hiscoreType, iterator, delims[len(delims)-1])
				}
			}

			c <- name.PageLimitInfo{Limit: start_page, HiscoreType: hiscoreType}
			<-concurrentGoroutines
		}(hiscoreSkills, hiscoreSkills.Limit, respChan, proxyIterator)

	}

	limits_set := 0

	for range name.HiscoreMinigames {
		r := <-respChan

		r.HiscoreType.Limit = r.Limit
		pageLimitManager.Add(r.HiscoreType)
		limits_set++

		secondsRan := time.Now().Unix() - start
		playersPerSecond := float64(limits_set) / float64(secondsRan)
		playersPerHour := (playersPerSecond * 3600) / 1000 // K players per hour.
		fmt.Printf("%d Limits found @ %.2fK/hr\n", limits_set, playersPerHour)

	}
	// wait for all routines to finish.
	wg.Wait()

	return pageLimitManager
}

func getValidPage(first_results map[string]bool, startPage int, delimPage int, hiscoreType name.HiscoreType, iterator *entities.ProxyIterator, finalDelimiter int) int {
	page := 1
	last_page := 0

	state := "continue"
	for page = startPage; page < 80_000; page += delimPage {
		results, _ := GetNames(hiscoreType, page, iterator)
		if mapEqual(first_results, results) {
			state = "finished"
			// print when finished all delimiters
			if delimPage == finalDelimiter {
				fmt.Printf("[%s] Checked %d [+=%d] (%s)\n", state, page, delimPage, hiscoreType.Minigames)
			}
			break
		} else {
			last_page = page // keep going.
			fmt.Printf("[%s] Checked %d [+=%d] (%s)\n", state, page, delimPage, hiscoreType.Minigames)
		}

	}
	return last_page
}

func mapEqual(m1 map[string]bool, m2 map[string]bool) bool {
	for k, v := range m1 {
		if m2[k] != v {
			return false
		}
	}
	return true
}