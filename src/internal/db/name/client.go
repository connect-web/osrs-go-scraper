package name

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/db"
)

type NameClient struct {
	*db.Client
}

func NewNameClient() *NameClient {
	return &NameClient{Client: db.NewDBClient()}
}

func (nameClient *NameClient) Close() {
	err := nameClient.Client.Close()
	if err != nil {
		fmt.Printf("Failed to close connection on nameClient.")
	}
}
