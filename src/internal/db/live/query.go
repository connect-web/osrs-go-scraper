package live

import (
	"encoding/json"
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utility/entities"
	"github.com/connect-web/Low-Latency/internal/utility/playerutils"
	"log"
	"strings"
)

// GET QUERIES

// returns a map of player name which require looking up on the Hiscores.
func (liveClient *LiveClient) fetchUnknownStats(limit int) (playerNameToId map[string]int, err error) {
	playerNameToId = make(map[string]int)
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
	-- ORDER BY players.first_seen // this order will slow things down.
	LIMIT ($1);
	`

	rows, err := liveClient.Client.DB.Query(selectQuery, limit)

	if err != nil {
		log.Println(err.Error())
		//log.Printf("Failed to get players with unknown stats.")
		return playerNameToId, err
	}

	defer rows.Close()

	for rows.Next() {
		var playerId int
		var name string

		scanErr := rows.Scan(&playerId, &name)
		if scanErr != nil {
			log.Printf("Failed to scan row from unknown stats query: %s\n", scanErr.Error())
		} else {
			playerNameToId[name] = playerId
		}
	}

	if err = rows.Err(); err != nil {
		log.Println(err) // just log it, if we have collected anything we can still return it.
	}

	return playerNameToId, nil
}

// POST QUERIES

func (liveClient *LiveClient) insertNotFound(playerIds map[int]struct{}) (err error) {
	if 65_000 < len(playerIds) {
		// batch the player_id inserts
		err = nil
		for _, playerIdBatch := range entities.ChunkIdsMap(playerIds, 50_000) {
			err = liveClient.insertNotFound(playerIdBatch)
		}
		return err
	}

	// Prepare usernames for insert transaction
	var valueStrings []string
	var valueArgs []interface{}
	i := 1
	for name := range playerIds {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d)", i))
		valueArgs = append(valueArgs, name)
		i++
	}

	// Create the insert query
	insertQuery := fmt.Sprintf(`
	INSERT INTO NOT_FOUND (PlayerId) 
	VALUES %s 
   	ON CONFLICT DO NOTHING`, strings.Join(valueStrings, ","))

	fmt.Printf("Inserting %d args\n", len(valueArgs))
	// Begin a transaction
	tx, err := liveClient.Client.DB.Begin()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return err
	}

	//fmt.Printf("Inserting: %d new name\n", len(playerIds))

	// Execute the insert query
	_, err = tx.Exec(insertQuery, valueArgs...)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			fmt.Println("Failed to rollback on NOT_FOUND Insert..?")
		}
		log.Println("Failed to bulk insert into Players:", err)
		return err
	}

	// Commit the transaction
	return tx.Commit()
}

func (liveClient *LiveClient) insertOrUpdatePlayers(simplePlayers []playerutils.SimplePlayer) (err error) {
	if 10_000 < len(simplePlayers) {
		// batch the player_id inserts
		err = nil
		for _, simplePlayerBatch := range entities.ChunkSimplePlayer(simplePlayers, 10_000) {
			err = liveClient.insertOrUpdatePlayers(simplePlayerBatch)
		}
		return err
	}

	var valueStrings []string
	var valueArgs []interface{}
	columnCount := 6
	i := 0
	for _, plr := range simplePlayers {
		if len(plr.Skills) == 0 && len(plr.Minigames) == 0 {
			// skip if skills & Minigames are empty.
			// this prevents errors where the plr object failed to load stats or API was down and returned empty json response.
			continue
		}
		skillsJson, skillErr := json.Marshal(plr.Skills)               // Convert map to JSON string
		skillLevelJson, skillLevelErr := json.Marshal(plr.SkillLevels) // Convert map to JSON string
		skillRatioJson, skillRatioErr := json.Marshal(plr.SkillRatios) // Convert map to JSON string
		minigames, minigameErr := json.Marshal(plr.Minigames)          // Convert map to JSON string

		// debug errors
		if skillErr != nil || minigameErr != nil || skillLevelErr != nil || skillRatioErr != nil {
			// Log each error if not nil
			if skillErr != nil {
				fmt.Println("skillErr:", skillErr.Error())
			}
			if minigameErr != nil {
				fmt.Printf("[%v] minigameErr: %s\n", plr.Minigames, minigameErr.Error())
			}
			if skillLevelErr != nil {
				fmt.Println("skillLevelErr:", skillLevelErr.Error())
			}
			if skillRatioErr != nil {
				fmt.Println("skillRatioErr:", skillRatioErr.Error())
			}
			continue // Skip this player if any error occurs
		}

		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", i*columnCount+1, i*columnCount+2, i*columnCount+3, i*columnCount+4, i*columnCount+5, i*columnCount+6))
		valueArgs = append(valueArgs, plr.PID, plr.LastUpdated, skillsJson, minigames, skillLevelJson, skillRatioJson)
		i++
	}
	if len(valueStrings) == 0 {
		return nil
	}

	// Create the insert query
	insertQuery := fmt.Sprintf(`
	INSERT INTO player_live (PlayerId, last_updated, skills_experience, minigames, skills_levels, skills_ratio) 
	VALUES %s 
   	ON CONFLICT (PlayerId) DO UPDATE SET
		skills_experience = EXCLUDED.skills_experience,
		minigames = EXCLUDED.minigames,
		skills_levels = EXCLUDED.skills_levels,
		skills_ratio = EXCLUDED.skills_ratio,
		last_updated = EXCLUDED.last_updated
	       `, strings.Join(valueStrings, ","))

	// Begin a transaction
	tx, err := liveClient.Client.DB.Begin()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return err
	}

	fmt.Printf("Inserting: %d new name\n", len(simplePlayers))

	// Execute the insert query
	_, err = tx.Exec(insertQuery, valueArgs...)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			fmt.Println("Failed to rollback on player_live Insert..?")
		}
		log.Println("Failed to bulk insert into player_live:", err)
		return err
	}

	// Commit the transaction
	return tx.Commit()
}
