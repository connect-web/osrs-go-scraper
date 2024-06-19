package main

import (
	"fmt"
	utils "github.com/connect-web/Low-Latency-Utils"
	"names.go/entities"
	"names.go/entities/limits"
	"sync"
	"time"
)

/*
The old school runescape hiscores has a limit of the top 2 million users per hiscore ranking

Some hiscore rankings do not have 2 million users participating in the activity, therefore this finds the page ranges that the Leaderboard nameFinder will visit
- This prevents wasted requests on pages that do not exist.

*/

func FindLimits() *limits.PageLimits {
	pageLimitManager := limits.NewPageLimitManager()

	concurrentGoroutines := make(chan struct{}, *maxNbConcurrentGoroutines)
	respChan := make(chan limits.PageLimitInfo)
	var wg sync.WaitGroup

	for _, hiscoreSkills := range limits.HiscoreMinigames {
		wg.Add(1)

		go func(hiscoreType limits.HiscoreType, page_limit int, c chan limits.PageLimitInfo, iterator *utils.ProxyIterator) {
			defer wg.Done()
			concurrentGoroutines <- struct{}{}
			start_page := 10

			first_page, _ := entities.GetNames(hiscoreType, 1, iterator)
			if len(first_page) == 0 {
				fmt.Printf("First page failed. %s\n", hiscoreType.Minigames)
			} else {
				delims := []int{5000, 1000, 100, 10}

				for _, delim_page := range delims {
					start_page = getValidPage(first_page, start_page, delim_page, hiscoreType, iterator, delims[len(delims)-1])
				}
			}

			c <- limits.PageLimitInfo{Limit: start_page, HiscoreType: hiscoreType}
			<-concurrentGoroutines
		}(hiscoreSkills, hiscoreSkills.Limit, respChan, proxyIterator)

	}

	limits_set := 0

	for range limits.HiscoreMinigames {
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

func getValidPage(first_results map[string]bool, startPage int, delimPage int, hiscoreType limits.HiscoreType, iterator *utils.ProxyIterator, finalDelimiter int) int {
	page := 1
	last_page := 0

	state := "continue"
	for page = startPage; page < 80_000; page += delimPage {
		results, _ := entities.GetNames(hiscoreType, page, iterator)
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
