package live

import (
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utility/entities"
	"github.com/connect-web/Low-Latency/internal/utility/playerutils"
	"log"
	"strconv"
	"strings"
)

func (liveClient *LiveClient) InsertOrUpdateLiveStats(players []playerutils.PlayerTotals) (err error) {
	if 10_000 < len(players) {
		// batch the player_id inserts
		err = nil
		for _, simplePlayerBatch := range entities.ChunkPlayer(players, 10_000) {
			err = liveClient.InsertOrUpdateLiveStats(simplePlayerBatch)
		}
		return err
	}

	var valueStrings []string
	var valueArgs []interface{}

	paramPrint := map[int]int{}

	column_count := 5
	i := 0
	for _, plr := range players {
		if plr.CombatLevel == 3 && plr.OverallExperience == 0 && plr.TotalLevel == 23 {
			// skip if skills & Minigames are empty.
			continue
		}
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*column_count+1, i*column_count+2, i*column_count+3, i*column_count+4, i*column_count+5))
		valueArgs = append(valueArgs, plr.PlayerId, plr.LastUpdated, plr.CombatLevel, plr.OverallExperience, plr.TotalLevel)
		paramPrint[plr.PlayerId] = plr.PlayerId
		i++
	}
	if len(valueStrings) == 0 {
		return nil
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
	tx, err := liveClient.Client.DB.Begin()
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
		if strings.Contains(err.Error(), "data type of parameter $") {
			paramNumber := strings.Split(err.Error(), "data type of parameter $")[1]
			fmt.Printf("param: %s.\n", paramNumber)
			integer, integerErr := strconv.Atoi(paramNumber)
			if integerErr == nil {
				fmt.Println(paramPrint[integer])
			}

		}
		return err
	}

	// Commit the transaction
	return tx.Commit()
}
