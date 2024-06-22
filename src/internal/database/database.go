package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

// DBClient represents a client to the database with a connection status.
type DBClient struct {
	DB              *sql.DB
	Connected       bool
	ConnectionError error
}

// NewDBClient initializes a new database client.
func NewDBClient() *DBClient {
	return &DBClient{
		Connected: false,
	}
}

// Connect establishes a connection to the database.
func (client *DBClient) Connect() error {
	user := os.Getenv("lowLatencyUser")
	password := os.Getenv("lowLatencyPassword")

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
func (client *DBClient) Close() error {
	if client.DB != nil {
		return client.DB.Close()
	}
	return nil
}

func test_connection() {
	dbClient := NewDBClient()
	err := dbClient.Connect()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer dbClient.Close()

	if dbClient.Connected {
		fmt.Println("Connected to database successfully")
		// Database operations go here
	}
}
