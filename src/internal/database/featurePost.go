package database

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utils/entities"
	"github.com/connect-web/Low-Latency/internal/utils/stats"
	"log"
	"strings"
)

func (uc *LeaderboardClient) InsertOrUpdateFeatures(players []stats.AdvancedPlayer) (err error) {
	if !uc.DBClient.Connected {
		log.Fatal("Not connected to database")
		// todo fatal might not be the best choice, reconnect would be a better solution.
	}

	if 10_000 < len(players) {
		// batch the player_id inserts
		err = nil
		for _, simplePlayerBatch := range entities.ChunkAdvancedPlayer(players, 10_000) {
			err = uc.InsertOrUpdateFeatures(simplePlayerBatch)
		}
		return err
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}
	column_count := 5
	for i, plr := range players {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*column_count+1, i*column_count+2, i*column_count+3, i*column_count+4, i*column_count+5))
		valueArgs = append(valueArgs, plr.PlayerId, plr.LastUpdated, plr.CombatLevel, plr.OverallExperience, plr.TotalLevel)
	}

	// Create the insert query
	insertQuery := fmt.Sprintf(`
	INSERT INTO player_live_stats (PlayerId, LAST_UPDATED, combat_level, Overall, total_level) 
	VALUES %s 
   	ON CONFLICT (PlayerId) DO UPDATE SET
		combat_level = EXCLUDED.combat_level,
		Overall = EXCLUDED.Overall,
		total_level = EXCLUDED.total_level,
		LAST_UPDATED = EXCLUDED.LAST_UPDATED
	       `, strings.Join(valueStrings, ","))

	// Begin a transaction
	tx, err := uc.DBClient.DB.Begin()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return err
	}

	fmt.Printf("Inserting: %d advanced players\n", len(players))

	// Execute the insert query
	_, err = tx.Exec(insertQuery, valueArgs...)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			fmt.Println("Failed to rollback on advanced players Insert/update..?")
		}
		log.Println("Failed to bulk insert into player_live_stats:", err)
		return err
	}

	// Commit the transaction
	return tx.Commit()
}
