package csv

import (
	"fmt"
	"strings"
)

// Preset represents a predefined CSV format
type Preset struct {
	Name        string
	Description string
	Mapping     ColumnMapping
}

// ColumnMapping defines how CSV columns map to BLEF fields
type ColumnMapping struct {
	BookID        string
	ISBN13        string
	ISBN10        string
	Title         string
	Author        string
	Language      string
	Publisher     string
	PublishedDate string
	Pages         string
	Rating        string
	Review        string
	Status        string
	DateRead      string
	DateAdded     string
	Tags          string
	Shelf         string
}

var (
	// GoodreadsPreset maps Goodreads CSV export format
	GoodreadsPreset = Preset{
		Name:        "goodreads",
		Description: "Goodreads library export",
		Mapping: ColumnMapping{
			// Don't map BookID - let the code use ISBN13 as ID
			ISBN13:        "ISBN13",
			ISBN10:        "ISBN",
			Title:         "Title",
			Author:        "Author",
			Publisher:     "Publisher",
			PublishedDate: "Year Published",
			Pages:         "Number of Pages",
			Rating:        "My Rating",
			Review:        "My Review",
			DateRead:      "Date Read",
			DateAdded:     "Date Added",
			Shelf:         "Exclusive Shelf",
		},
	}

	// BabelioPreset maps Babelio CSV export format
	BabelioPreset = Preset{
		Name:        "babelio",
		Description: "Babelio library export",
		Mapping: ColumnMapping{
			ISBN13:        "EAN",
			Title:         "Titre",
			Author:        "Auteur",
			Publisher:     "Editeur",
			PublishedDate: "Date de publication",
			Rating:        "Note",
			Review:        "Critique",
			Status:        "État",
			Shelf:         "Étagère",
		},
	}

	AllPresets = []Preset{GoodreadsPreset, BabelioPreset}
)

// DetectPreset attempts to automatically detect the CSV format
func DetectPreset(data *CSVData) *Preset {
	// Try Goodreads detection
	if hasColumns(data, []string{"Book Id", "Title", "Author", "ISBN13", "My Rating"}) {
		return &GoodreadsPreset
	}

	// Try Babelio detection
	if hasColumns(data, []string{"EAN", "Titre", "Auteur", "Étagère"}) {
		return &BabelioPreset
	}

	return nil
}

// hasColumns checks if all specified columns exist in the CSV
func hasColumns(data *CSVData, columns []string) bool {
	for _, col := range columns {
		if data.GetColumnIndex(col) < 0 {
			return false
		}
	}
	return true
}

// MapStatus converts platform-specific status to BLEF status
func MapStatus(value string, preset *Preset) string {
	if preset == nil {
		return normalizeStatus(value)
	}

	value = strings.TrimSpace(strings.ToLower(value))

	switch preset.Name {
	case "goodreads":
		switch value {
		case "read":
			return "read"
		case "currently-reading":
			return "reading"
		case "to-read":
			return "to-read"
		default:
			return "to-read"
		}

	case "babelio":
		switch value {
		case "lu", "read":
			return "read"
		case "en cours", "reading":
			return "reading"
		case "à lire", "to-read":
			return "to-read"
		case "abandonné", "abandoned":
			return "abandoned"
		default:
			return "to-read"
		}
	}

	return normalizeStatus(value)
}

// normalizeStatus attempts to normalize any status string
func normalizeStatus(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))

	// Common mappings
	if strings.Contains(value, "read") && !strings.Contains(value, "reading") {
		return "read"
	}
	if strings.Contains(value, "reading") || strings.Contains(value, "current") {
		return "reading"
	}
	if strings.Contains(value, "to") && strings.Contains(value, "read") {
		return "to-read"
	}
	if strings.Contains(value, "abandon") {
		return "abandoned"
	}
	if strings.Contains(value, "wish") {
		return "wishlist"
	}

	// Default
	return "to-read"
}

// CleanGoodreadsValue removes Excel formulas from Goodreads exports
// Converts ="value" or =""value"" to value
func CleanGoodreadsValue(value string) string {
	value = strings.TrimSpace(value)

	// Remove Excel formula wrapper: ="..." -> ...
	if strings.HasPrefix(value, `="`) && strings.HasSuffix(value, `"`) {
		value = strings.TrimPrefix(value, `="`)
		value = strings.TrimSuffix(value, `"`)

		// Remove internal double quotes: ""123"" -> 123
		value = strings.Trim(value, `"`)
	}

	return strings.TrimSpace(value)
}

// MapRating converts rating strings to float64
func MapRating(value string) float64 {
	value = strings.TrimSpace(value)
	if value == "" || value == "0" {
		return 0
	}

	// Try parsing as number
	var rating float64
	if _, err := fmt.Sscanf(value, "%f", &rating); err == nil {
		if rating > 5 {
			rating = rating / 2 // Convert 10-point scale to 5-point
		}
		if rating < 0 {
			rating = 0
		}
		if rating > 5 {
			rating = 5
		}
		return rating
	}

	return 0
}
