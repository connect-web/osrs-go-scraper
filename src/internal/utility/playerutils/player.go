package playerutils

import "time"

type PlayerTotals struct {
	PlayerId          int
	TotalLevel        int16
	CombatLevel       int16
	OverallExperience int64
	LastUpdated       time.Time
}

func NewPlayerTotals() PlayerTotals {
	return PlayerTotals{
		PlayerId:          0,
		TotalLevel:        0,
		CombatLevel:       0,
		OverallExperience: 0,
	}
}

type Ratios struct {
	PlayerId    int
	LastUpdated time.Time
}

/*
	- Possible Ratios -

- Could order all players by X Skill DESC
- Then could compare the skill ratios for players around it but for all skills looking at pearson
- can set a threshold for pearson e.g:  [0.9]



Store all players in a list and then get the index and get all players + 50 index and - 50 index if exists.
- then run the pearson check.

Start off doing it for 1 skill and set the filter to no minigames.

Test skill : Cooking

*/

func (ap *PlayerTotals) Calculate(sp SimplePlayer) {
	ap.OverallExperience = sp.calculateOverallExperience()
	ap.CombatLevel = sp.calculateCombat()
	ap.TotalLevel = sp.calculateOverallLevel()
	ap.LastUpdated = sp.LastUpdated
	ap.PlayerId = sp.PID
}
