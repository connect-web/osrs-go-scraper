package nameutils

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type HtmlHiscore struct {
	DevPrints       bool
	SkillIndexes    map[int]string
	MinigameIndexes map[int]string
	AvoidRows       []string
	page            int
}

func NewHtmlHiscore(page int) *HtmlHiscore {
	return &HtmlHiscore{
		page: page,

		DevPrints: false,
		SkillIndexes: map[int]string{
			1: "Name",
			2: "Rank",
			3: "Level",
			4: "Experience",
		},
		MinigameIndexes: map[int]string{
			1: "Name",
			2: "Rank",
			3: "Score",
		},
		AvoidRows: []string{
			"Personal scores",
			"XP",
		},
	}
}

func (df *HtmlHiscore) ExtendedStringToInteger(text string) int {
	text = strings.ReplaceAll(text, ",", "")
	val, _ := strconv.Atoi(text)
	return val
}

func (df *HtmlHiscore) FilterRow(tds *goquery.Selection, minigame bool) map[string]string {
	var titleMap map[int]string
	if minigame {
		titleMap = df.MinigameIndexes
	} else {
		titleMap = df.SkillIndexes
	}

	rowData := make(map[string]string)
	tds.Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}
		for _, avoidText := range df.AvoidRows {
			if strings.Contains(s.Text(), avoidText) {
				return
			}
		}

		title := titleMap[i]
		if title != "" {
			rowData[title] = s.Text()
		} else {
			//fmt.Printf("[Error] Couldn't find title for: %s\n", s.Text())
		}
	})

	return rowData
}

func (df *HtmlHiscore) MapValueToInteger(d map[string]string) map[string]int {
	newDict := make(map[string]int)
	for k, v := range d {
		newDict[k] = df.ExtendedStringToInteger(v)
	}
	return newDict
}

func (df *HtmlHiscore) GetUsernamesFromHtmlString(htmlPage string) map[string]bool {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlPage))
	if err != nil {
		panic(err)
	}

	Users := make(map[string]bool)

	doc.Find("tr").Each(func(i int, row *goquery.Selection) {

		minigameRow := false
		if href, exists := row.Find("a").Attr("href"); exists && strings.Contains(href, "hiscorepersonal") {
			minigameRow = true
		}

		tds := row.Find("td")

		rowData := df.FilterRow(tds, minigameRow)

		if minigameRow {

			username, validName := rowData["Name"]
			minigameRank, validRank := rowData["Rank"]

			if validName && validRank {
				username = strings.ReplaceAll(username, "\n", "")
				minigameRank = strings.ReplaceAll(minigameRank, "\n", "")
				minigameRank = strings.ReplaceAll(minigameRank, ",", "")

				username = SanitizeString(username)
				// debug print
				//fmt.Printf("[Page %d] %s %s\n", df.page, username, minigameRank)
				Users[username] = true
			}
		}
	})

	return Users
}

// SanitizeString replaces HTML entities and invalid UTF-8 sequences
func SanitizeString(input string) string {
	// Decode HTML entities
	decoded := decodeHTMLEntities(input)

	// Replace non-breaking space with a regular space
	decoded = strings.ReplaceAll(decoded, "\u00A0", " ")

	// Remove invalid UTF-8 characters
	sanitized := removeInvalidUTF8(decoded)

	return sanitized
}

// decodeHTMLEntities decodes HTML entities in a string
func decodeHTMLEntities(s string) string {
	// Create a reader for the HTML-encoded string
	r := strings.NewReader(s)

	// Parse the HTML and extract text
	var buf bytes.Buffer
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return buf.String()
		case html.TextToken:
			buf.Write(z.Text())
		}
	}
}

// removeInvalidUTF8 removes invalid UTF-8 characters from a string
func removeInvalidUTF8(s string) string {
	// Use the charmap package to replace invalid UTF-8 characters
	encoder := charmap.ISO8859_1.NewDecoder()
	result, _, _ := transform.String(encoder, s)

	// Validate UTF-8 again, remove invalid sequences if any
	var validUTF8 bytes.Buffer
	for i, r := range result {
		if r == utf8.RuneError {
			_, size := utf8.DecodeRuneInString(result[i:])
			if size == 1 {
				// Skip invalid byte
				continue
			}
		}
		validUTF8.WriteRune(r)
	}

	return validUTF8.String()
}

// debugging html:
func writePage(content string) {
	file, err := os.Create("page.html")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
