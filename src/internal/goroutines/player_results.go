package goroutines

import "github.com/connect-web/Low-Latency/internal/utility/playerutils"

type PlayerLookupResults struct {
	Players      []playerutils.SimplePlayer
	PlayerTotals []playerutils.PlayerTotals
	NotFound     map[int]struct{}
}
