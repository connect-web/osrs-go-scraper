package nameFilter

import (
	"fmt"
	"regexp"
)

func Filter(usernames map[string]bool) map[string]bool {
	filtered_usernames := map[string]bool{}

	for username := range usernames {
		filtered_username := replaceUnicodeSpaces(
			replaceUnicodeNbsp(
				replaceUnicodeDash(username)))
		filtered_usernames[filtered_username] = true
	}
	return filtered_usernames
}

func replaceUnicodeSpaces(input string) string {
	// Directly include the non-breaking space character
	nbsp := "\\xao"                // Non-breaking space as a Unicode code point
	re := regexp.MustCompile(nbsp) // Compile a regex with the non-breaking space
	// Replace matched non-breaking space with a regular space
	return re.ReplaceAllString(input, " ")
}

func replaceUnicodeNbsp(input string) string {
	// Directly include the non-breaking space character
	nbsp := "\u00A0"               // Non-breaking space as a Unicode code point
	re := regexp.MustCompile(nbsp) // Compile a regex with the non-breaking space
	// Replace matched non-breaking space with a regular space
	return re.ReplaceAllString(input, " ")
}

func replaceUnicodeDash(input string) string {
	// Directly include the non-breaking space character
	dash := "\\u2011"              // Non-breaking space as a Unicode code point
	re := regexp.MustCompile(dash) // Compile a regex with the non-breaking space
	// Replace matched non-breaking space with a regular space
	return re.ReplaceAllString(input, "-")
}

func main() {
	input := "Galar tank"
	fmt.Println("Original:", input)
	fmt.Println("Code points of original name:")
	printCodePoints(input)
	output := replaceUnicodeSpaces(input)
	fmt.Println("Modified:", output)
	fmt.Println("Code points of modified name:")
	printCodePoints(output)
}

// printCodePoints prints the Unicode code points of each character in the string
func printCodePoints(name string) {
	for _, runeValue := range name {
		fmt.Printf("%U ", runeValue)
	}
	fmt.Println()
}
