package stats

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"time"
)

type JsonResponse struct {
	Skills []struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Rank  int    `json:"rank"`
		Level int    `json:"level"`
		Xp    int    `json:"xp"`
	} `json:"skills"`
	Activities []struct {
		Id    int    `json:"id"`
		Name  string `json:"name"`
		Rank  int    `json:"rank"`
		Score int    `json:"score"`
	} `json:"activities"`
}

type SimplePlayer struct {
	Username    string
	LastUpdated time.Time
	PID         int
	Skills      map[string]int // Skill : Experience directly from Hiscores
	Minigames   map[string]int // Minigame/activity : score directly from Hiscores
	SkillLevels map[string]int
	SkillRatios map[string]float32 // 32 bits will have enough useful data
}

func NewSimplePlayer(PID int) *SimplePlayer {
	return &SimplePlayer{
		PID:         PID,
		Skills:      make(map[string]int),
		Minigames:   make(map[string]int),
		SkillLevels: make(map[string]int),
		SkillRatios: make(map[string]float32),
	}
}

func (sp *SimplePlayer) LoadJson(response JsonResponse) bool {
	valid := false
	for _, val := range response.Skills {
		if val.Xp != 0 {
			valid = true
		}
		if 0 < val.Xp {
			sp.Skills[val.Name] = val.Xp
		}
	}
	for _, val := range response.Activities {
		if 0 < val.Score {
			sp.Minigames[val.Name] = val.Score
		}
		if val.Score != 0 {
			valid = true
		}
	}
	return valid

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

func (sp *SimplePlayer) FromSQL(rows *sql.Rows) error {
	var skillsData, minigamesData []byte

	err := rows.Scan(&sp.PID, &sp.Username, &skillsData, &minigamesData)
	if err != nil {
		return err
	}

	err = json.Unmarshal(skillsData, &sp.Skills)
	if err != nil {
		return err
	}

	err = json.Unmarshal(minigamesData, &sp.Minigames)
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
