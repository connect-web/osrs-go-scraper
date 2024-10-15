package gain

import (
	"encoding/json"
	"fmt"
	"github.com/connect-web/Low-Latency/internal/utility/entities"
	"github.com/connect-web/Low-Latency/internal/utility/playerutils"
	"log"
	"strings"
)

func (client *GainsClient) fetchPlayersRequireGains(limit int) (outdatedPlayers []playerutils.SimplePlayer, err error) {
	outdatedPlayers = make([]playerutils.SimplePlayer, 0)

	selectQuery := `
	SELECT
		players.id,
		players.name,
		pl.last_updated,
		pl.skills_experience,
		pl.minigames
	FROM players
	LEFT JOIN not_found NF ON NF.playerid = players.id
	LEFT JOIN player_live PL ON PL.playerid = players.id
	LEFT JOIN player_gains pg ON pg.playerid = players.id
	where
		nf.playerid is null 
	  AND PL.last_updated is not null
	  AND NOW() - pl.last_updated  > INTERVAL '3 day'
	  AND (
	      PG.last_updated is null OR NOW() - PG.last_updated > INTERVAL '7 DAY'
	  )
	LIMIT $1;
	`

	rows, err := client.Client.DB.Query(selectQuery, limit)

	if err != nil {
		log.Println(err.Error())
		return outdatedPlayers, err
	}

	defer rows.Close()

	for rows.Next() {
		player := playerutils.SimplePlayer{
			Skills:    map[string]int{},
			Minigames: map[string]int{},
		}
		var skillBytes, minigameBytes []byte

		scanErr := rows.Scan(&player.PID, &player.Username, &player.LastUpdated, &skillBytes, &minigameBytes)
		if scanErr != nil {
			log.Printf("Failed to scan row from unknown stats query: %s\n", scanErr.Error())
			continue
		}
		if unmarshalLevelErr := json.Unmarshal(skillBytes, &player.Skills); unmarshalLevelErr != nil {
			log.Println(unmarshalLevelErr.Error())
			continue
		}
		if unmarshalExperienceErr := json.Unmarshal(minigameBytes, &player.Minigames); unmarshalExperienceErr != nil {
			log.Println(unmarshalExperienceErr.Error())
			continue
		}
		outdatedPlayers = append(outdatedPlayers, player)

	}

	if err = rows.Err(); err != nil {
		log.Println(err) // just log it, if we have collected anything we can still return it.
	}

	return outdatedPlayers, nil
}

func (client *GainsClient) InsertOrUpdateGains(simplePlayers []playerutils.SimplePlayer) (err error) {
	if 10_000 < len(simplePlayers) {
		// batch the player_id inserts
		err = nil
		for _, simplePlayerBatch := range entities.ChunkSimplePlayer(simplePlayers, 10_000) {
			err = client.InsertOrUpdateGains(simplePlayerBatch)
		}
		return err
	}

	var valueStrings []string
	var valueArgs []interface{}
	columnCount := 4
	i := 0
	for _, plr := range simplePlayers {
		// Convert map to JSON string
		skillsJson, skillErr := json.Marshal(plr.Skills)
		skillRatioJson, skillRatioErr := json.Marshal(plr.SkillRatios)
		minigamesJson, minigameErr := json.Marshal(plr.MinigamesDaily)

		// debug errors
		if skillErr != nil || minigameErr != nil || skillRatioErr != nil {
			// Log each error if not nil
			if skillErr != nil {
				fmt.Println("skillErr:", skillErr.Error())
			}
			if minigameErr != nil {
				fmt.Printf("[%v] minigameErr: %s\n", plr.MinigamesDaily, minigameErr.Error())
			}
			if skillRatioErr != nil {
				fmt.Println("skillRatioErr:", skillRatioErr.Error())
			}
			continue // Skip this player if any error occurs
		}

		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*columnCount+1, i*columnCount+2, i*columnCount+3, i*columnCount+4))
		valueArgs = append(valueArgs, plr.PID, skillsJson, minigamesJson, skillRatioJson)
		i++
	}

	// Create the insert query
	insertQuery := fmt.Sprintf(`
	INSERT INTO player_gains (PlayerId, skills_experience, minigames, skills_ratio) 
	VALUES %s 
   	ON CONFLICT (PlayerId) DO UPDATE SET
		skills_experience = EXCLUDED.skills_experience,
		minigames = EXCLUDED.minigames,
		skills_ratio = EXCLUDED.skills_ratio,
		last_updated = NOW()
	       `, strings.Join(valueStrings, ","))

	// Begin a transaction
	tx, err := client.Client.DB.Begin()
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
