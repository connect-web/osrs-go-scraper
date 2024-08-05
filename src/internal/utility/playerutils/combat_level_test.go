package playerutils

import "testing"

func TestCombatLevel(t *testing.T) {
	_ = CombatLevel(42, 40, 40, 41, 44, 35, 41)
}
