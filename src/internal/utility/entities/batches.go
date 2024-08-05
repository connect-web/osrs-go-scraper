package entities

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utility/playerutils"
)

// todo replace all with interface

func ChunkNamesMap(names map[string]struct{}, chunkSize int) []map[string]struct{} {
	var chunks []map[string]struct{}
	temp := map[string]struct{}{}
	for name := range names {
		temp[name] = struct{}{}
		if len(temp) == chunkSize {
			chunks = append(chunks, temp)
			temp = map[string]struct{}{}
		}
	}
	if 0 < len(temp) {
		chunks = append(chunks, temp)
		temp = map[string]struct{}{}
	}
	fmt.Printf("Turned %d name into %d lists of %d name!\n",
		len(names), len(chunks), chunkSize,
	)
	return chunks
}

func ChunkNames(names []string, chunkSize int) [][]string {
	var chunks [][]string
	var temp []string
	for _, name := range names {
		temp = append(temp, name)
		if len(temp) == chunkSize {
			chunks = append(chunks, temp)
			temp = []string{}
		}
	}
	if 0 < len(temp) {
		chunks = append(chunks, temp)
		temp = []string{}
	}
	fmt.Printf("Turned %d name into %d lists of %d name!\n",
		len(names), len(chunks), chunkSize,
	)
	return chunks
}

func ChunkIdsMap(ids map[int]struct{}, chunkSize int) []map[int]struct{} {
	var chunks []map[int]struct{}
	temp := map[int]struct{}{}
	for name := range ids {
		temp[name] = struct{}{}
		if len(temp) == chunkSize {
			chunks = append(chunks, temp)
			temp = map[int]struct{}{}
		}
	}
	if 0 < len(temp) {
		chunks = append(chunks, temp)
		temp = map[int]struct{}{}
	}
	fmt.Printf("Turned %d ids into %d lists of %d ids!\n",
		len(ids), len(chunks), chunkSize,
	)
	return chunks
}

func ChunkIds(ids []int, chunkSize int) [][]int {
	var chunks [][]int
	var temp []int
	for _, name := range ids {
		temp = append(temp, name)
		if len(temp) == chunkSize {
			chunks = append(chunks, temp)
			temp = []int{}
		}
	}
	if 0 < len(temp) {
		chunks = append(chunks, temp)
		temp = []int{}
	}
	fmt.Printf("Turned %d ids into %d lists of %d ids!\n",
		len(ids), len(chunks), chunkSize,
	)
	return chunks
}

func ChunkSimplePlayer(ids []playerutils.SimplePlayer, chunkSize int) [][]playerutils.SimplePlayer {
	var chunks [][]playerutils.SimplePlayer
	var temp []playerutils.SimplePlayer
	for _, name := range ids {
		temp = append(temp, name)
		if len(temp) == chunkSize {
			chunks = append(chunks, temp)
			temp = []playerutils.SimplePlayer{}
		}
	}
	if 0 < len(temp) {
		chunks = append(chunks, temp)
		temp = []playerutils.SimplePlayer{}
	}
	fmt.Printf("Turned %d SimplePlayers into %d lists of %d SimplePlayers!\n",
		len(ids), len(chunks), chunkSize,
	)
	return chunks
}

func ChunkPlayer(ids []playerutils.PlayerTotals, chunkSize int) [][]playerutils.PlayerTotals {
	var chunks [][]playerutils.PlayerTotals
	var temp []playerutils.PlayerTotals
	for _, name := range ids {
		temp = append(temp, name)
		if len(temp) == chunkSize {
			chunks = append(chunks, temp)
			temp = []playerutils.PlayerTotals{}
		}
	}
	if 0 < len(temp) {
		chunks = append(chunks, temp)
		temp = []playerutils.PlayerTotals{}
	}
	fmt.Printf("Turned %d Players into %d lists of %d Players!\n",
		len(ids), len(chunks), chunkSize,
	)
	return chunks
}

func ChunkUserMap(userMap map[string]int, chunkSize int) []map[string]int {
	var chunks []map[string]int
	temp := map[string]int{}
	for username, playerId := range userMap {

		temp[username] = playerId

		if len(temp) == chunkSize {
			chunks = append(chunks, temp)
			temp = map[string]int{}
		}
	}
	if 0 < len(temp) {
		chunks = append(chunks, temp)
		temp = map[string]int{}
	}
	fmt.Printf("Turned %d users into %d lists of %d users!\n",
		len(userMap), len(chunks), chunkSize,
	)
	return chunks
}

func ChunkPearsonResults(data []playerutils.PearsonResults, chunkSize int) [][]playerutils.PearsonResults {
	var chunks [][]playerutils.PearsonResults
	var temp []playerutils.PearsonResults

	for _, item := range data {
		temp = append(temp, item)
		if len(temp) == chunkSize {
			chunks = append(chunks, temp)
			temp = []playerutils.PearsonResults{}
		}
	}
	if 0 < len(temp) {
		chunks = append(chunks, temp)
		temp = []playerutils.PearsonResults{}
	}
	fmt.Printf("Turned %d Pearson results into %d lists of %d results!\n",
		len(data), len(chunks), chunkSize,
	)
	return chunks
}
