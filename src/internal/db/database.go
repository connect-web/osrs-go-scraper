package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"time"
)

// Client represents a client to the database with a connection status.
type Client struct {
	DB              *sql.DB
	Connected       bool
	ConnectionError error
}

// NewDBClient initializes a new database client.
func NewDBClient() *Client {
	client := &Client{
		Connected: false,
	}
	err := client.Connect()
	if err != nil {
		fmt.Printf("failed to connect to db: %s\n", err.Error())
		var retries int
		for !client.Connected {
			retries++
			fmt.Printf("%d: Retrying db connection...", retries)
			_ = client.Connect()
			time.Sleep(30 * time.Second)
		}
	}
	return client
}

// Connect establishes a connection to the database.
func (client *Client) Connect() error {
	/*



		user := os.Getenv("lowLatencyWebUser")
			password := os.Getenv("lowLatencyWebPassword")
			host := os.Getenv("lowLatencyWebHost")
			port := os.Getenv("lowLatencyWebPort")
			dbname := os.Getenv("lowLatencyWebDatabase")

	*/

	user := os.Getenv("lowLatencyUser")

	password := os.Getenv("lowLatencyPassword")
	//host := os.Getenv("lowLatencyDevHost")
	host := os.Getenv("lowLatencyHost")

	port := os.Getenv("lowLatencyPort")
	dbname := os.Getenv("lowLatencyDatabase")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		client.ConnectionError = err
		return err
	}

	// Try to make a connection
	err = db.Ping()
	if err != nil {
		client.ConnectionError = err
		return err
	}

	client.DB = db
	client.Connected = true
	client.ConnectionError = nil
	return nil
}

// Close terminates the connection to the database.
func (client *Client) Close() error {
	if client.DB != nil {
		return client.DB.Close()
	}
	return nil
}
