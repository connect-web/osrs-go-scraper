package playerutils

import (
	"fmt"
	"testing"
)

func TestLevelCalculator(t *testing.T) {
	calculator := NewLevelCalculator()
	level, err := calculator.GetLevel(13_000_000)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("Level: %d\n", level)
}
