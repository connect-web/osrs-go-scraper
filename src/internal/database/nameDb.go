package database

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utils/entities"
	"github.com/connect-web/Low-Latency/internal/utils/name"
	"github.com/lib/pq"
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

func SubmitUsernames(usernames map[string]struct{}) error {
	filtered_usernames := name.Filter(usernames)
	err := connect_and_submit(filtered_usernames)
	return err
}

func connect_and_submit(usernames map[string]struct{}) error {
	dbClient := NewDBClient()
	err := dbClient.Connect()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return err
	}
	defer func(dbClient *DBClient) {
		dbClientErr := dbClient.Close()
		if dbClientErr != nil {
			fmt.Println(dbClientErr.Error())
		}
	}(dbClient)

	Client := NewLeaderboardClient(dbClient)

	// first make sure the usernames are not in old_players
	validUsernames, validNameErr := Client.getValidUsernames(usernames)
	if validNameErr != nil {
		fmt.Printf("Failed to validate usernames %s\n", validNameErr.Error())
		return validNameErr
	}
	fmt.Printf("Removed %d old players from username list!\n", len(usernames)-len(validUsernames))

	// insert the usernames
	err = Client.insertUsernames(validUsernames)
	if err == nil {
		fmt.Println("Successfully saved usernames.")
	}

	return err
}

func (uc *LeaderboardClient) getValidUsernames(usernames map[string]struct{}) (validUsernames map[string]struct{}, err error) {
	validUsernames = make(map[string]struct{})
	existingNames := make(map[string]struct{})

	if !uc.DBClient.Connected {
		log.Fatal("Not connected to database")
		// todo fatal might not be the best choice, reconnect would be a better solution.
	}

	// Create the query
	// find all name that exist in old_players already from username map
	selectQuery := `
		SELECT
		name
	FROM old_players
	where
		name = any(($1));
	`

	usernameArray := []string{}
	for username := range usernames {
		usernameArray = append(usernameArray, username)
	}

	rows, err := uc.DBClient.DB.Query(selectQuery, pq.Array(usernameArray))

	if err != nil {
		log.Println(err.Error())
		log.Printf("Failed to get players with unknown stats.")
		return map[string]struct{}{}, err
	}

	defer rows.Close()

	// fill existingNames from the name that exist in old_players that match.
	for rows.Next() {
		var name string
		scanErr := rows.Scan(&name)
		if scanErr != nil {
			log.Printf("Failed to scan row from unknown stats query: %s\n", scanErr.Error())
		} else {
			existingNames[name] = struct{}{}
		}
	}

	if err = rows.Err(); err != nil {
		log.Println(err) // just log it, if we have collected anything we can still return it.
	}

	// add all the usernames from original map into valid if they do not exist in old_names database.
	for username := range usernames {
		_, exists := existingNames[username]
		if !exists {
			validUsernames[username] = struct{}{}
		}
	}

	return validUsernames, nil
}

// GetUser retrieves a user by ID.
func (uc *LeaderboardClient) insertUsernames(names map[string]struct{}) error {
	if 30_000 < len(names) {
		for _, batch := range entities.ChunkNamesMap(names, 10_000) {
			_ = uc.insertUsernames(batch)
		}
		return nil

	}

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

	fmt.Printf("Inserting: %d new name\n", len(names))

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
