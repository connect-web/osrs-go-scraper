package database

import (
	"log"
)

// returns a map of player name which require looking up on the Hiscores.
func (uc *LeaderboardClient) fetchUnknownStats(limit int) (playerNameToId map[string]int, err error) {
	playerNameToId = make(map[string]int)
	if !uc.DBClient.Connected {
		log.Fatal("Not connected to database")
		// todo fatal might not be the best choice, reconnect would be a better solution.
	}

	// Create the query
	// Find all players where they do not exist in not_found table && do not exist in player_live table.
	// ordered by First seen
	selectQuery := `
		SELECT
		players.id,
		players.name
	FROM players
	LEFT JOIN not_found NF ON NF.playerid = players.id
	LEFT JOIN player_live PL ON PL.playerid = players.id
	where
		nf.playerid is null
	AND PL.playerid is null
	ORDER BY players.first_seen
	LIMIT ($1);
	`

	rows, err := uc.DBClient.DB.Query(selectQuery, limit)

	if err != nil {
		log.Println(err.Error())
		//log.Printf("Failed to get players with unknown stats.")
		return playerNameToId, err
	}

	defer rows.Close()

	for rows.Next() {
		var player_id int
		var name string

		scanErr := rows.Scan(&player_id, &name)
		if scanErr != nil {
			log.Printf("Failed to scan row from unknown stats query: %s\n", scanErr.Error())
		} else {
			playerNameToId[name] = player_id
		}
	}

	if err = rows.Err(); err != nil {
		log.Println(err) // just log it, if we have collected anything we can still return it.
	}

	return playerNameToId, nil
}

func (uc *LeaderboardClient) fetchOutdatedStats(limit int) (playerNameToId map[string]int, err error) {
	playerNameToId = make(map[string]int)
	if !uc.DBClient.Connected {
		log.Fatal("Not connected to database")
		// todo fatal might not be the best choice, reconnect would be a better solution.
	}

	// Create the query
	// Find all players where they do not exist in not_found table && do not exist in player_live table.
	// ordered by First seen
	selectQuery := `
		SELECT
		players.id,
		players.name
	FROM players
	LEFT JOIN not_found NF ON NF.playerid = players.id
	LEFT JOIN player_live PL ON PL.playerid = players.id
	where
		nf.playerid is null
	AND PL.playerid is not null
	AND PL.LAST_UPDATED  <= NOW() - INTERVAL '7 days' -- If the player_live was last updated over 7 days ago.
	ORDER BY PL.LAST_UPDATED
	LIMIT ($1);
	`

	rows, err := uc.DBClient.DB.Query(selectQuery, limit)

	if err != nil {
		log.Println(err.Error())
		//log.Printf("Failed to get players with unknown stats.")
		return playerNameToId, err
	}

	defer rows.Close()

	for rows.Next() {
		var player_id int
		var name string

		scanErr := rows.Scan(&player_id, &name)
		if scanErr != nil {
			log.Printf("Failed to scan row from unknown stats query: %s\n", scanErr.Error())
		} else {
			playerNameToId[name] = player_id
		}
	}

	if err = rows.Err(); err != nil {
		log.Println(err) // just log it, if we have collected anything we can still return it.
	}

	return playerNameToId, nil
}
