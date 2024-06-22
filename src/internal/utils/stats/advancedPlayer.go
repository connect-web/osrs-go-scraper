package stats

import "time"

type AdvancedPlayer struct {
	PlayerId          int
	TotalLevel        int16
	CombatLevel       int16
	OverallExperience int64
	LastUpdated       time.Time
}

func (ap *AdvancedPlayer) Calculate(sp SimplePlayer) {
	ap.OverallExperience = sp.calculateOverallExperience()
	ap.CombatLevel = sp.calculateCombat()
	ap.TotalLevel = sp.calculateOverallLevel()
	ap.LastUpdated = sp.LastUpdated
	ap.PlayerId = sp.PID
}
