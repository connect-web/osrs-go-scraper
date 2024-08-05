package live

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/db"
)

type LiveClient struct {
	*db.Client
}

func NewLiveClient() *LiveClient {
	return &LiveClient{Client: db.NewDBClient()}
}

func (liveClient *LiveClient) Close() {
	err := liveClient.Client.Close()
	if err != nil {
		fmt.Printf("Failed to close connection on liveClient.")
	}
}
