package main

import (
	"fmt"
	"log"
	"strings"
)

// UserClient extends DBClient with specific methods for handling users.
type LeaderboardClient struct {
	*DBClient
}

// NewUserClient creates a new UserClient.
func NewLeaderboardClient(dbClient *DBClient) *LeaderboardClient {
	return &LeaderboardClient{DBClient: dbClient}
}

// GetUser retrieves a user by ID.
func (uc *LeaderboardClient) insertUsernames(names map[string]bool) error {
	if !uc.DBClient.Connected {
		log.Fatal("Not connected to database")
	}

	// Prepare usernames for insert transaction
	valueStrings := []string{}
	valueArgs := []interface{}{}
	i := 1
	for name := range names {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d)", i))
		valueArgs = append(valueArgs, name)
		i++
	}

	// Create the insert query
	insertQuery := fmt.Sprintf(`
	INSERT INTO Players (Name) 
	VALUES %s 
   	ON CONFLICT DO NOTHING`, strings.Join(valueStrings, ","))

	// Begin a transaction
	tx, err := uc.DBClient.DB.Begin()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return err
	}

	fmt.Printf("Inserting: %d new names\n", len(names))

	// Execute the insert query
	_, err = tx.Exec(insertQuery, valueArgs...)
	if err != nil {
		tx.Rollback()
		log.Println("Failed to bulk insert into Players:", err)
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

func main() {
	dbClient := NewDBClient()
	err := dbClient.Connect()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer dbClient.Close()

	Client := NewLeaderboardClient(dbClient)
	Client.DBClient.Connect()

	usernames := map[string]bool{
		"test123":     true,
		"test12343":   true,
		"test123213":  true,
		"test1235534": true,
	}

	err = Client.insertUsernames(usernames) // Example usage
	if err == nil {
		fmt.Println("Successfully saved usernames.")
	} else {
		log.Fatal(err.Error())
	}
}
