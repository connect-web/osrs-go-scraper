package metrics

import "time"

var (
	start             = time.Now().Unix()
	StatsFound        = 0
	PlayersInserted   = 0
	NotfoundInserted  = 0
	PlayerTotals      = 0
	UsernamesInserted = 0
)

func GetHourly(size int) float64 {
	secondsRan := time.Now().Unix() - start
	PerSecond := float64(size) / float64(secondsRan)
	PerHour := (PerSecond * 3600) / 1000 // K players per hour.
	return PerHour
}
