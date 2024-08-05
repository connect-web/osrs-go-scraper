package playerutils

import (
	"math"
)

func (sp *SimplePlayer) calculateCombat() int16 {
	combatLvl := CombatLevel(sp.SkillLevels["Attack"], sp.SkillLevels["Defence"], sp.SkillLevels["Strength"], sp.SkillLevels["Hitpoints"], sp.SkillLevels["Prayer"], sp.SkillLevels["Ranged"], sp.SkillLevels["Magic"])
	return combatLvl
}

func minimumCombatLevel(combatLevel int16) int16 {
	var minimumLevel int16 = 3
	if combatLevel < minimumLevel {
		return minimumLevel
	}
	return combatLevel
}

func CombatLevel(attack, defence, strength, hitpoints, prayer, ranged, magic int) int16 {
	attack = int(removeVirtualLevel(int16(attack)))
	defence = int(removeVirtualLevel(int16(defence)))
	strength = int(removeVirtualLevel(int16(strength)))
	hitpoints = int(removeVirtualLevel(int16(hitpoints)))
	prayer = int(removeVirtualLevel(int16(prayer)))
	ranged = int(removeVirtualLevel(int16(ranged)))
	magic = int(removeVirtualLevel(int16(magic)))

	base := float64(defence+hitpoints+prayer/2) * 0.25

	melee := float64(attack+strength) * 0.325
	rng := float64(math.Floor(float64(ranged)*1.5)) * 0.325
	mage := float64(math.Floor(float64(magic)*1.5)) * 0.325
	maximumLevel := math.Max(melee, math.Max(rng, mage))
	level := math.Round((base+maximumLevel)*1000) / 1000
	return minimumCombatLevel(int16(math.Floor(level)))
}
