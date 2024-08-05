package gain

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/db"
)

type GainsClient struct {
	*db.Client
}

func NewGainsClient() *GainsClient {
	return &GainsClient{Client: db.NewDBClient()}
}

func (client *GainsClient) Close() {
	err := client.Client.Close()
	if err != nil {
		fmt.Printf("Failed to close connection on client.")
	}
}
