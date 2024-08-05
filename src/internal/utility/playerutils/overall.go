package playerutils

var totalSkills = 23

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

	// Since the map does not store level 1's it will need to be added to the total level.
	skillsTrained := len(sp.SkillLevels)
	untrainedSkills := totalSkills - skillsTrained

	totalLevel += int16(untrainedSkills)

	return totalLevel
}

func (sp *SimplePlayer) calculateOverallExperience() int64 {
	var totalExperience int64
	for _, experience := range sp.Skills {
		totalExperience = totalExperience + int64(experience)
	}
	return totalExperience
}
