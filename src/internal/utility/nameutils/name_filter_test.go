package nameutils

import (
	"fmt"
	"testing"
)

func TestNameFilter(t *testing.T) {
	input := "Galar tank"
	fmt.Println("Original:", input)
	fmt.Println("Code points of original name:")
	printCodePoints(input)
	output := replaceUnicodeNbsp(input)
	fmt.Println("Modified:", output)
	fmt.Println("Code points of modified name:")
	printCodePoints(output)
}
