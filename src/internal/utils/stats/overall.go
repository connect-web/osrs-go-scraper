package stats

func removeVirtualLevel(level int16) int16 {
	var maxLevel int16 = 99
	if maxLevel < level {
		return maxLevel
	}
	return level
}

func (sp *SimplePlayer) calculateOverallLevel() int16 {
	var totalLevel int16
	for _, level := range sp.SkillLevels {
		actualLevel := removeVirtualLevel(int16(level))
		totalLevel = totalLevel + actualLevel
	}
	return totalLevel
}

func (sp *SimplePlayer) calculateOverallExperience() int64 {
	var totalExperience int64
	for _, experience := range sp.Skills {
		totalExperience = totalExperience + int64(experience)
	}
	return totalExperience
}
