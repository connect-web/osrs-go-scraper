package playerutils

import (
	"time"
)

func (sp *SimplePlayer) CalculateGains(onlinePlayer SimplePlayer) (playerGains SimplePlayer) {
	playerGains = SimplePlayer{
		Username:       sp.Username,
		LastUpdated:    sp.LastUpdated,
		PID:            sp.PID,
		Skills:         make(map[string]int),
		MinigamesDaily: make(map[string]float32),
		SkillRatios:    make(map[string]float32),
	}

	secondsDuration := time.Now().Unix() - sp.LastUpdated.Unix()
	daysDuration := float64(secondsDuration) / 86400
	// The gains will be Mean gains per day

	for skillName, skillExperience := range onlinePlayer.Skills {
		oldExperience := sp.Skills[skillName]
		gainedExperience := skillExperience - oldExperience
		if 0 < gainedExperience {
			//fmt.Printf("%s [%d] gained %d experience in %s\n", onlinePlayer.Username, onlinePlayer.PID, gainedExperience, skillName)
			meanDailyExperience := int(float64(gainedExperience) / daysDuration)
			playerGains.Skills[skillName] = meanDailyExperience
		}
	}

	for minigameName, minigameScore := range onlinePlayer.Minigames {
		oldScore := sp.Minigames[minigameName]
		gainedScore := minigameScore - oldScore
		if 0 < gainedScore {
			//fmt.Printf("%s [%d] gained %d experience in %s\n", onlinePlayer.Username, onlinePlayer.PID, gainedScore, minigameName)
			meanDailyScore := float32(float64(gainedScore) / daysDuration)
			playerGains.MinigamesDaily[minigameName] = meanDailyScore
		}
	}

	_ = playerGains.calculateRatios()
	// error just means 0 skill experience but does not matter, just store null in place of skill Level map.
	return playerGains
}
