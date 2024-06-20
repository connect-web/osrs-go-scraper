package entities

import "errors"

func (sp *SimplePlayer) Calculations() {
	_ = sp.calculateRatios()
	_ = sp.calculateLevels()
}

func (sp *SimplePlayer) calculateRatios() error {
	// First calculate the overall Experience
	// If the player is not in the top 2 Million Overall it will not be displayed therefore we will calculate it
	// Some will have this but to be uniform we will calculate all.
	var overall_experience int64
	for name, experience := range sp.Skills {
		if name == "Overall" {
			continue
		}
		overall_experience += int64(experience)
	}

	if overall_experience == 0 {
		// cannot divide by 0 so skills_ratio will be null
		return errors.New("Cannot divide by zero.")
	}

	for name, experience := range sp.Skills {
		if name == "Overall" {
			continue
		}
		//                     100 million      /   3 billion
		ratio_precision := float64(experience) / float64(overall_experience)
		// ratio precision is now between 0 and 1
		stored_ratio := float32(ratio_precision)
		// compact float 32 for storage but the precise division has been completed.
		sp.SkillRatios[name] = stored_ratio
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
		level, err := calculator.GetLevel(experience)
		if err == nil {
			sp.SkillLevels[name] = level
		}

	}
	return nil
}
