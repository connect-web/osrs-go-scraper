package name

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/metrics"
	"github.com/connect-web/Low-Latency/internal/requests/leaderboards"
	"github.com/connect-web/Low-Latency/internal/utility/entities"
	"github.com/connect-web/Low-Latency/internal/utility/nameutils"
	"sync"
)

/*
The old school Runescape hiscores has a limit of the top 2 million users per hiscore ranking

Some hiscore rankings do not have 2 million users participating in the activity, therefore this finds the page ranges that the Leaderboard nameFinder will visit
- This prevents wasted requests on pages that do not exist.

*/

func FindLimits(maxNbConcurrentGoroutines *int) *nameutils.PageLimits {
	pageLimitManager := nameutils.NewPageLimitManager()
	concurrentGoroutines := make(chan struct{}, *maxNbConcurrentGoroutines)
	respChan := make(chan nameutils.PageLimitInfo)
	var wg sync.WaitGroup

	for _, hiscoreSkills := range nameutils.HiscoreMinigames {
		wg.Add(1)

		go func(hiscoreType nameutils.HiscoreType, pageLimit int, c chan nameutils.PageLimitInfo, iterator *entities.ProxyIterator) {
			defer wg.Done()
			concurrentGoroutines <- struct{}{}
			startPage := 10

			firstPage, _ := leaderboards.GetNames(hiscoreType, 1, iterator)
			if len(firstPage) == 0 {
				fmt.Printf("First page failed. %s\n", hiscoreType.Minigames)
			} else {
				delimiters := []int{5000, 1000, 100, 10}

				for _, delimPage := range delimiters {
					startPage = getValidPage(firstPage, startPage, delimPage, hiscoreType, iterator, delimiters[len(delimiters)-1])
				}
			}

			c <- nameutils.PageLimitInfo{Limit: startPage, HiscoreType: hiscoreType}
			<-concurrentGoroutines
		}(hiscoreSkills, hiscoreSkills.Limit, respChan, entities.ProxyList)

	}

	limitsSet := 0

	for range nameutils.HiscoreMinigames {
		r := <-respChan

		r.HiscoreType.Limit = r.Limit
		pageLimitManager.Add(r.HiscoreType)
		limitsSet++
		fmt.Printf("%d Limits found @ %.2fK/hr\n", limitsSet, metrics.GetHourly(limitsSet))

	}
	// wait for all routines to finish.
	wg.Wait()

	return pageLimitManager
}

func getValidPage(firstResults map[string]bool, startPage int, delimPage int, hiscoreType nameutils.HiscoreType, iterator *entities.ProxyIterator, finalDelimiter int) int {
	page := 1
	lastPage := 0

	state := "continue"
	for page = startPage; page < 80_000; page += delimPage {
		results, _ := leaderboards.GetNames(hiscoreType, page, iterator)
		if mapEqual(firstResults, results) {
			state = "finished"
			// print when finished all delimiters
			if delimPage == finalDelimiter {
				fmt.Printf("[%s] Checked %d [+=%d] (%s)\n", state, page, delimPage, hiscoreType.Minigames)
			}
			break
		} else {
			lastPage = page // keep going.
			fmt.Printf("[%s] Checked %d [+=%d] (%s)\n", state, page, delimPage, hiscoreType.Minigames)
		}

	}
	return lastPage
}

func mapEqual(m1 map[string]bool, m2 map[string]bool) bool {
	for k, v := range m1 {
		if m2[k] != v {
			return false
		}
	}
	return true
}
