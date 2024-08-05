package playerutils

import (
	"github.com/connect-web/Low-Latency/internal/utility/nameutils"
	statistics "github.com/montanaflynn/stats"
	"log"
	"time"
)

type PearsonResults struct {
	PlayerId        int
	LastUpdated     time.Time
	Skill           string
	LinkedPlayerIds []int
}

func (sp *SimplePlayer) CompareRatioVariance(otherPlayer SimplePlayer) float64 {
	// get the pearson across each skill
	// take the mean
	var mainPlayerRatios []float64
	var altPlayerRatios []float64

	for _, skill := range nameutils.HiscoreSkills {
		mainPlayerRatio := sp.SkillRatios[skill.Skill]
		alternatePlayerRatio := otherPlayer.SkillRatios[skill.Skill]

		mainPlayerRatios = append(mainPlayerRatios, float64(mainPlayerRatio))
		altPlayerRatios = append(altPlayerRatios, float64(alternatePlayerRatio))

	}
	currentUserRatio := statistics.LoadRawData(mainPlayerRatios)
	otherUserRatio := statistics.LoadRawData(altPlayerRatios)
	result, err := statistics.Pearson(currentUserRatio, otherUserRatio)
	if err == nil {
		//fmt.Printf("%.2f\n", result)
		return result
	} else {
		log.Printf("Pearson Err: %s from comparing:\n%v\n%v\n", err.Error(), mainPlayerRatios, altPlayerRatios)

	}
	return 0
}
