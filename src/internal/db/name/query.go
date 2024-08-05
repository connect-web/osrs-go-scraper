package name

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utility/entities"
	"github.com/lib/pq"
	"log"
	"strings"
)

func (nameClient *NameClient) getValidUsernames(usernames map[string]struct{}) (validUsernames map[string]struct{}, err error) {
	validUsernames = make(map[string]struct{})
	existingNames := make(map[string]struct{})

	// find all name that exist in old_players already from username map
	selectQuery := "SELECT name FROM players where name = any(($1));"

	var usernameArray []string
	for username := range usernames {
		usernameArray = append(usernameArray, username)
	}

	rows, err := nameClient.Client.DB.Query(selectQuery, pq.Array(usernameArray))

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

func (nameClient *NameClient) getAllUsernames() (err error) {
	rows, err := nameClient.Client.DB.Query("SELECT name FROM players")

	if err != nil {
		// must be fatal as this sort of error will impact performance massively.
		// this only runs on startup though.
		log.Fatalf("Failed to load usernames from Players: %s\n", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		scanErr := rows.Scan(&name)
		if scanErr != nil {
			log.Printf("Failed to scan row from unknown stats query: %s\n", scanErr.Error())
			continue
		}
		// loading names into memory.
		knownUsernames[name] = struct{}{}
	}
	return nil
}

func (nameClient *NameClient) insertUsernames(names map[string]struct{}) error {
	if 30_000 < len(names) {
		for _, batch := range entities.ChunkNamesMap(names, 10_000) {
			_ = nameClient.insertUsernames(batch)
		}
		return nil

	}

	if !nameClient.Client.Connected {
		log.Fatal("Not connected to database")
	}

	// Prepare usernames for insert transaction
	var valueStrings []string
	var valueArgs []interface{}
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
	tx, err := nameClient.Client.DB.Begin()
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

	commitErr := tx.Commit()

	if commitErr == nil && !lowMemory {
		// if it's not low memory mode then add the newly added usernames to the known_usernames.
		for name := range names {
			knownUsernames[name] = struct{}{}
		}
	}

	// Commit the transaction
	return commitErr
}
