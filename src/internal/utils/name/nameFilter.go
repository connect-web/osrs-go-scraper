package name

import (
	"fmt"
	"regexp"
	"strings"
)

func Filter(usernames map[string]struct{}) map[string]struct{} {
	filtered_usernames := map[string]struct{}{}

	for username := range usernames {
		filtered_username :=
			replaceUnicodeNbsp(
				replaceUnicodeDash(username))
		filtered_usernames[filtered_username] = struct{}{}
	}
	return filtered_usernames
}

func replaceUnicodeNbsp_old(input string) string {
	// Directly include the non-breaking space character
	nbsp := "\u00A0"               // Non-breaking space as a Unicode code point
	re := regexp.MustCompile(nbsp) // Compile a regex with the non-breaking space
	// Replace matched non-breaking space with a regular space

	return re.ReplaceAllString(input, " ")
}

func replaceUnicodeNbsp(input string) string {
	incorrectNbsp := string([]byte{0xA0}) // Simulating the incorrect single-byte input

	// Replace incorrect non-breaking space with a regular space
	output := strings.ReplaceAll(input, incorrectNbsp, " ")
	return output
}

func replaceUnicodeDash(input string) string {
	// Correctly include the non-breaking hyphen as a Unicode code point
	dash := `\x{2011}`             // Non-breaking hyphen in hexadecimal notation
	re := regexp.MustCompile(dash) // Compile a regex with the non-breaking hyphen
	// Replace matched non-breaking hyphen with a regular hyphen
	return re.ReplaceAllString(input, "-")
}

func test_name_filter() {
	input := "Galar tank"
	fmt.Println("Original:", input)
	fmt.Println("Code points of original name:")
	printCodePoints(input)
	output := replaceUnicodeNbsp(input)
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
