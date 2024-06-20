package main

import (
	"flag"
	"fmt"
	utils "github.com/connect-web/Low-Latency-Utils"
	"log"
	"stats.go/entities"
	"time"
)

var (
	start                     = time.Now().Unix()
	threads                   = 20
	maxNbConcurrentGoroutines = flag.Int("MaxRoutines", threads, "The number of goroutines that are allowed to run concurrently")
	proxyIterator             = utils.NewProxyIterator("proxies.txt")
	usernames_verified        = 0
)

func main() {
	username := "zezima"
	player, err := entities.Get_player_stats(username, 432, proxyIterator)

	if err == nil {
		fmt.Println(player)
	} else {
		log.Printf(err.Error())
	}
}
