package entities

import (
	"fmt"
)

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
	fmt.Printf("Turned %d names into %d lists of %d names!\n",
		len(names), len(chunks), chunkSize,
	)
	return chunks
}
