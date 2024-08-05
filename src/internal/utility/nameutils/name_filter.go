package nameutils

import (
	"fmt"
	"regexp"
	"strings"
)

func Filter(usernames map[string]struct{}) map[string]struct{} {
	filteredUsernames := map[string]struct{}{}

	for username := range usernames {
		filteredUsername := FilterName(username)
		filteredUsernames[filteredUsername] = struct{}{}
	}
	return filteredUsernames
}

func FilterName(username string) string {
	return replaceUnicodeNbsp(
		replaceUnicodeDash(username))
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

// printCodePoints prints the Unicode code points of each character in the string
func printCodePoints(name string) {
	for _, runeValue := range name {
		fmt.Printf("%U ", runeValue)
	}
	fmt.Println()
}
