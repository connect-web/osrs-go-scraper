package database

import (
	"encoding/json"
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utils/entities"
	"github.com/connect-web/Low-Latency/internal/utils/stats"

	"log"
	"strings"
)

func (uc *LeaderboardClient) insertNotFound(playerIds map[int]struct{}) (err error) {
	if !uc.DBClient.Connected {
		log.Fatal("Not connected to database")
		// todo fatal might not be the best choice, reconnect would be a better solution.
	}

	if 65_000 < len(playerIds) {
		// batch the player_id inserts
		err = nil
		for _, playerIdBatch := range entities.ChunkIdsMap(playerIds, 50_000) {
			err = uc.insertNotFound(playerIdBatch)
		}
		return err
	}

	// Prepare usernames for insert transaction
	valueStrings := []string{}
	valueArgs := []interface{}{}
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

	// Begin a transaction
	tx, err := uc.DBClient.DB.Begin()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return err
	}

	fmt.Printf("Inserting: %d new name\n", len(playerIds))

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

func (uc *LeaderboardClient) insertOrUpdatePlayers(simplePlayers []stats.SimplePlayer) (err error) {
	if !uc.DBClient.Connected {
		log.Fatal("Not connected to database")
		// todo fatal might not be the best choice, reconnect would be a better solution.
	}

	if 10_000 < len(simplePlayers) {
		// batch the player_id inserts
		err = nil
		for _, simplePlayerBatch := range entities.ChunkSimplePlayer(simplePlayers, 10_000) {
			err = uc.insertOrUpdatePlayers(simplePlayerBatch)
		}
		return err
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}
	column_count := 5
	for i, plr := range simplePlayers {
		skillsJson, skillErr := json.Marshal(plr.Skills)          // Convert map to JSON string
		minigamesJson, minigameErr := json.Marshal(plr.Minigames) // Convert map to JSON string

		skillLevelJson, skillLevelErr := json.Marshal(plr.SkillLevels) // Convert map to JSON string
		skillRatioJson, skillRatioErr := json.Marshal(plr.SkillRatios) // Convert map to JSON string

		// debug errors
		if skillErr != nil || minigameErr != nil || skillLevelErr != nil || skillRatioErr != nil {
			// Log each error if not nil
			if skillErr != nil {
				fmt.Println("skillErr:", skillErr.Error())
			}
			if minigameErr != nil {
				fmt.Println("minigameErr:", minigameErr.Error())
			}
			if skillLevelErr != nil {
				fmt.Println("skillLevelErr:", skillLevelErr.Error())
			}
			if skillRatioErr != nil {
				fmt.Println("skillRatioErr:", skillRatioErr.Error())
			}
			continue // Skip this player if any error occurs
		}

		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*column_count+1, i*column_count+2, i*column_count+3, i*column_count+4, i*column_count+5))
		valueArgs = append(valueArgs, plr.PID, skillsJson, minigamesJson, skillLevelJson, skillRatioJson)
	}

	// Create the insert query
	insertQuery := fmt.Sprintf(`
	INSERT INTO player_live (PlayerId, skills_experience, minigames, skills_levels, skills_ratio) 
	VALUES %s 
   	ON CONFLICT (PlayerId) DO UPDATE SET
		skills_experience = EXCLUDED.skills_experience,
		minigames = EXCLUDED.minigames,
		skills_levels = EXCLUDED.skills_levels,
		skills_ratio = EXCLUDED.skills_ratio,
		last_updated = NOW()
	       `, strings.Join(valueStrings, ","))

	// Begin a transaction
	tx, err := uc.DBClient.DB.Begin()
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
			fmt.Println("Failed to rollback on NOT_FOUND Insert..?")
		}
		log.Println("Failed to bulk insert into player_live:", err)
		return err
	}

	// Commit the transaction
	return tx.Commit()
}
