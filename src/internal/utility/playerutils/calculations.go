package playerutils

import "errors"

func (sp *SimplePlayer) Calculations() {
	_ = sp.calculateRatios()
	_ = sp.calculateLevels()
}

func (sp *SimplePlayer) calculateRatios() error {
	if sp.SkillRatios == nil {
		sp.SkillRatios = make(map[string]float32)
	}
	// First calculate the overall Experience
	// If the player is not in the top 2 Million Overall it will not be displayed therefore we will calculate it
	// Some will have this but to be uniform we will calculate all.
	var overallExperience int64
	for name, experience := range sp.Skills {
		if name == "Overall" || experience <= 0 {
			continue
		}
		overallExperience += int64(experience)

	}

	if overallExperience < 0 {
		// cannot divide by 0 so skills_ratio will be null
		return errors.New("cannot divide by zero")
	}
	for name, experience := range sp.Skills {
		if name == "Overall" || experience <= 0 {
			continue
		}
		//                     100 million      /   3 billion
		ratioPrecision := float64(experience) / float64(overallExperience)
		// ratio precision is now between 0 and 1
		storedRatio := float32(ratioPrecision)
		// compact float 32 for storage but the precise division has been completed.
		sp.SkillRatios[name] = storedRatio
	}

	return nil

}

func (sp *SimplePlayer) calculateLevels() error {
	// Converts your experience into Runescape Levels including virtual levels.
	calculator := NewLevelCalculator()
	for name, experience := range sp.Skills {
		if name == "Overall" {
			continue
		}
		if experience <= 0 {
			// no experience is level 1.
			sp.SkillLevels[name] = 1
		}
		level, err := calculator.GetLevel(experience)
		if err == nil {
			sp.SkillLevels[name] = level
		}

	}
	return nil
}
