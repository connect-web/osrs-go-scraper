package leaderboards

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
)

type SimplePlayer struct {
	Username  string
	Skills    map[string]int
	Minigames map[string]int
	PID       int
}

func NewSimplePlayer() *SimplePlayer {
	return &SimplePlayer{
		Skills:    make(map[string]int),
		Minigames: make(map[string]int),
	}
}

func (sp *SimplePlayer) Add(earnedSkills map[string]int, earnedMinigames map[string]int) {
	for k, v := range earnedSkills {
		sp.Skills[k] = v
	}

	for k, v := range earnedMinigames {
		sp.Minigames[k] = v
	}
}

func (sp *SimplePlayer) String() string {
	return fmt.Sprintf("Player (%s)\nSkills:\n%v,\nMinigames:\n%v\n", sp.Username, sp.Skills, sp.Minigames)
}

func (sp *SimplePlayer) Dump() map[string]interface{} {
	return map[string]interface{}{
		"Username":  sp.Username,
		"Skills":    sp.Skills,
		"Minigames": sp.Minigames,
		"PID":       sp.PID,
	}
}

func (player *SimplePlayer) FromSQL(rows *sql.Rows) error {
	var skillsData, minigamesData []byte

	err := rows.Scan(&player.PID, &player.Username, &skillsData, &minigamesData)
	if err != nil {
		return err
	}

	err = json.Unmarshal(skillsData, &player.Skills)
	if err != nil {
		return err
	}

	err = json.Unmarshal(minigamesData, &player.Minigames)
	return err
}

func (sp *SimplePlayer) Compare(other *SimplePlayer) map[string]map[string]int {
	skillDiff := make(map[string]int)
	minigameDiff := make(map[string]int)

	for skill, value := range sp.Skills {

		otherValue, ok := other.Skills[skill]

		if !ok {
			otherValue = 0
		}
		diff := value - otherValue
		if diff < 0 {
			diff = -diff
		}
		//fmt.Printf("[SKILL] %s , %d VS %d\n", skill, value, otherValue)
		//fmt.Println(diff)
		if diff != 0 {

			skillDiff[skill] = int(math.Max(float64(value), float64(otherValue)))
		}
	}
	//fmt.Println(skillDiff)

	for minigame, value := range sp.Minigames {
		otherValue, ok := other.Minigames[minigame]

		if !ok {
			otherValue = 0
		}
		diff := value - otherValue
		if diff < 0 {
			diff = -diff
		}
		if diff != 0 {
			minigameDiff[minigame] = int(math.Max(float64(value), float64(otherValue)))
		}
	}

	return map[string]map[string]int{
		"Skills":    skillDiff,
		"Minigames": minigameDiff,
	}
}

func IsEmptyDiff(both map[string]map[string]int) bool {
	return len(both["Skills"]) == 0 && len(both["Minigames"]) == 0
}
