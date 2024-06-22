package entities

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utils/stats"
)

// todo replace all with interface

func ChunkNamesMap(names map[string]struct{}, chunkSize int) []map[string]struct{} {
	chunks := []map[string]struct{}{}
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
	chunks := [][]string{}
	temp := []string{}
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
	chunks := []map[int]struct{}{}
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
	chunks := [][]int{}
	temp := []int{}
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

func ChunkSimplePlayer(ids []stats.SimplePlayer, chunkSize int) [][]stats.SimplePlayer {
	chunks := [][]stats.SimplePlayer{}
	temp := []stats.SimplePlayer{}
	for _, name := range ids {
		temp = append(temp, name)
		if len(temp) == chunkSize {
			chunks = append(chunks, temp)
			temp = []stats.SimplePlayer{}
		}
	}
	if 0 < len(temp) {
		chunks = append(chunks, temp)
		temp = []stats.SimplePlayer{}
	}
	fmt.Printf("Turned %d SimplePlayers into %d lists of %d SimplePlayers!\n",
		len(ids), len(chunks), chunkSize,
	)
	return chunks
}

func ChunkAdvancedPlayer(ids []stats.AdvancedPlayer, chunkSize int) [][]stats.AdvancedPlayer {
	chunks := [][]stats.AdvancedPlayer{}
	temp := []stats.AdvancedPlayer{}
	for _, name := range ids {
		temp = append(temp, name)
		if len(temp) == chunkSize {
			chunks = append(chunks, temp)
			temp = []stats.AdvancedPlayer{}
		}
	}
	if 0 < len(temp) {
		chunks = append(chunks, temp)
		temp = []stats.AdvancedPlayer{}
	}
	fmt.Printf("Turned %d AdvancedPlayers into %d lists of %d AdvancedPlayers!\n",
		len(ids), len(chunks), chunkSize,
	)
	return chunks
}

func ChunkUserMap(userMap map[string]int, chunkSize int) []map[string]int {
	chunks := []map[string]int{}
	temp := map[string]int{}
	for username, player_id := range userMap {

		temp[username] = player_id

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
