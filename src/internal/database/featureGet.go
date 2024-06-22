package database

import (
	"encoding/json"
	"github.com/connect-web/Low-Latency/internal/utils/stats"
	"log"
)

// returns a map of player name which require looking up on the Hiscores.
func (uc *LeaderboardClient) FetchUnknownFeatures(limit int) (players []stats.SimplePlayer, err error) {
	players = make([]stats.SimplePlayer, 0)
	if !uc.DBClient.Connected {
		log.Fatal("Not connected to database")
		// todo fatal might not be the best choice, reconnect would be a better solution.
	}

	selectQuery := `
	SELECT
		pl.playerid,
		pl.last_updated,
		pl.skills_levels,
		pl.skills_experience
	FROM player_live pl
	LEFT JOIN player_live_stats pls on pls.playerid = pl.playerid
	WHERE
		pls.playerid is null
	limit ($1);
	`

	rows, err := uc.DBClient.DB.Query(selectQuery, limit)

	if err != nil {
		log.Println(err.Error())
		return players, err
	}

	defer rows.Close()

	for rows.Next() {
		var skillLevel, skillExperience []byte
		var player = stats.SimplePlayer{}

		scanErr := rows.Scan(&player.PID, &player.LastUpdated, &skillLevel, &skillExperience)
		if scanErr != nil {
			log.Printf("Failed to scan row from unknown stats query: %s\n", scanErr.Error())
			continue
		}
		if unmarshalLevelErr := json.Unmarshal(skillLevel, &player.SkillLevels); unmarshalLevelErr != nil {
			log.Println(unmarshalLevelErr.Error())
			continue
		}
		if unmarshalExperienceErr := json.Unmarshal(skillExperience, &player.Skills); unmarshalExperienceErr != nil {
			log.Println(unmarshalExperienceErr.Error())
			continue
		}

		players = append(players, player)
	}

	if err = rows.Err(); err != nil {
		log.Println(err) // just log it, if we have collected anything we can still return it.
	}

	return players, nil
}

func (uc *LeaderboardClient) FetchOutdatedFeatures(limit int) (players []stats.SimplePlayer, err error) {
	players = make([]stats.SimplePlayer, 0)
	if !uc.DBClient.Connected {
		log.Fatal("Not connected to database")
		// todo fatal might not be the best choice, reconnect would be a better solution.
	}

	selectQuery := `
	SELECT
		pl.playerid,
		pl.last_updated,
		pl.skills_levels,
		pl.skills_experience
	FROM player_live pl
			 LEFT JOIN player_live_stats pls on pls.playerid = pl.playerid
	WHERE
		pls.playerid is not null
		AND pl.last_updated != pls.last_updated
	limit ($1);
	`

	rows, err := uc.DBClient.DB.Query(selectQuery, limit)

	if err != nil {
		log.Println(err.Error())
		//log.Printf("Failed to get players with unknown stats.")
		return players, err
	}

	defer rows.Close()

	for rows.Next() {
		var skillLevel, skillExperience []byte
		var player = stats.SimplePlayer{}

		scanErr := rows.Scan(&player.PID, &player.LastUpdated, &skillLevel, &skillExperience)
		if scanErr != nil {
			log.Printf("Failed to scan row from unknown stats query: %s\n", scanErr.Error())
			continue
		}
		if unmarshalLevelErr := json.Unmarshal(skillLevel, &player.SkillLevels); unmarshalLevelErr != nil {
			log.Println(unmarshalLevelErr.Error())
			continue
		}
		if unmarshalExperienceErr := json.Unmarshal(skillExperience, &player.Skills); unmarshalExperienceErr != nil {
			log.Println(unmarshalExperienceErr.Error())
			continue
		}

		players = append(players, player)
	}

	if err = rows.Err(); err != nil {
		log.Println(err) // just log it, if we have collected anything we can still return it.
	}

	return players, nil
}
