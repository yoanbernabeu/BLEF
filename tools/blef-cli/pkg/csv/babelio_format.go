package csv

import (
	"fmt"
	"strings"

	"github.com/yoanbernabeu/BLEF/tools/blef-cli/pkg/blef"
)

// BabelioFormat implements CSVFormat for Babelio library exports
type BabelioFormat struct{}

func (f *BabelioFormat) Name() string {
	return "babelio"
}

func (f *BabelioFormat) Description() string {
	return "Babelio library export"
}

func (f *BabelioFormat) Detect(data *CSVData) bool {
	// Check for Babelio-specific columns
	// Real Babelio exports use "ISBN", "Titre", "Auteur", "Statut"
	requiredColumns := []string{"ISBN", "Titre", "Auteur", "Statut"}
	for _, col := range requiredColumns {
		if data.GetColumnIndex(col) < 0 {
			return false
		}
	}
	return true
}

func (f *BabelioFormat) GetImportMapping() ColumnMapping {
	return ColumnMapping{
		ISBN13:        "ISBN", // Real Babelio uses "ISBN" not "EAN"
		Title:         "Titre",
		Author:        "Auteur",
		Publisher:     "Editeur",
		PublishedDate: "Date de publication",
		Rating:        "Note",
		Status:        "Statut", // Real Babelio uses "Statut" not "État"
		DateAdded:     "Date d`entrée dans Babelio",
		// Note: Real Babelio exports don't have "Critique" or "Étagère" columns
	}
}

func (f *BabelioFormat) CleanValue(value string) string {
	// Babelio doesn't use special formatting, just trim
	return strings.TrimSpace(value)
}

func (f *BabelioFormat) MapStatus(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))

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

func (f *BabelioFormat) MapRating(value string) float64 {
	value = strings.TrimSpace(value)
	if value == "" || value == "0" {
		return 0
	}

	var rating float64
	if _, err := fmt.Sscanf(value, "%f", &rating); err == nil {
		// Babelio uses 5-point scale
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

func (f *BabelioFormat) GetExportHeaders() []string {
	// Real Babelio export format
	return []string{
		"ISBN", // Real format uses "ISBN" not "EAN"
		"Titre",
		"Auteur",
		"Editeur",
		"Date de publication",
		"Date d`entrée dans Babelio",
		"Statut", // Real format uses "Statut" not "État"
		"Note",
	}
}

func (f *BabelioFormat) ExportBook(book *blef.Book, entry *blef.Entry) []string {
	row := make([]string, len(f.GetExportHeaders()))

	// ISBN (ISBN13 or use book ID if not an ISBN)
	if book.Identifiers.ISBN13 != "" {
		row[0] = book.Identifiers.ISBN13
	} else {
		row[0] = book.ID // Fallback to book ID (might be custom ID like "SI19412249693")
	}

	// Titre
	row[1] = book.Title

	// Auteur
	if len(book.Authors) > 0 {
		authorNames := make([]string, len(book.Authors))
		for i, author := range book.Authors {
			authorNames[i] = author.Name
		}
		row[2] = strings.Join(authorNames, ", ")
	}

	// Editeur
	if book.Edition != nil {
		row[3] = book.Edition.Publisher
		// Date de publication
		row[4] = book.Edition.PublishedDate
	}

	// Entry data
	if entry != nil {
		// Date d'entrée dans Babelio
		if entry.UserData.AddedAt != nil {
			row[5] = entry.UserData.AddedAt.Format("2006-01-02 15:04:05")
		}

		// Statut
		row[6] = mapStatusToBabelio(entry.UserData.Status)

		// Note - export even if 0
		row[7] = fmt.Sprintf("%.1f", entry.UserData.Rating)
	}

	return row
}

// mapStatusToBabelio converts BLEF status to Babelio status
// Uses the exact format from real Babelio exports (with capital first letter)
func mapStatusToBabelio(status string) string {
	switch status {
	case "read":
		return "Lu"
	case "reading":
		return "En cours"
	case "to-read":
		return "A lire" // Note: Babelio uses "A" without accent
	case "abandoned":
		return "Abandonné"
	default:
		return "A lire"
	}
}
