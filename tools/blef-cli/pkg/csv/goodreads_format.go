package csv

import (
	"fmt"
	"strings"
	"time"

	"github.com/yoanbernabeu/BLEF/tools/blef-cli/pkg/blef"
)

// GoodreadsFormat implements CSVFormat for Goodreads library exports
type GoodreadsFormat struct{}

func (f *GoodreadsFormat) Name() string {
	return "goodreads"
}

func (f *GoodreadsFormat) Description() string {
	return "Goodreads library export"
}

func (f *GoodreadsFormat) Detect(data *CSVData) bool {
	// Check for Goodreads-specific columns
	requiredColumns := []string{"Book Id", "Title", "Author", "ISBN13", "My Rating"}
	for _, col := range requiredColumns {
		if data.GetColumnIndex(col) < 0 {
			return false
		}
	}
	return true
}

func (f *GoodreadsFormat) GetImportMapping() ColumnMapping {
	return ColumnMapping{
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
		Status:        "Exclusive Shelf", // In Goodreads, shelf = status
		DateRead:      "Date Read",
		DateAdded:     "Date Added",
		Shelf:         "Exclusive Shelf",
	}
}

func (f *GoodreadsFormat) CleanValue(value string) string {
	// Goodreads exports use Excel formulas like ="value" or =""value""
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

func (f *GoodreadsFormat) MapStatus(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))

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
}

func (f *GoodreadsFormat) MapRating(value string) float64 {
	value = strings.TrimSpace(value)
	if value == "" || value == "0" {
		return 0
	}

	var rating float64
	if _, err := fmt.Sscanf(value, "%f", &rating); err == nil {
		// Goodreads uses 5-point scale
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

func (f *GoodreadsFormat) GetExportHeaders() []string {
	return []string{
		"Book Id",
		"Title",
		"Author",
		"Author l-f",
		"Additional Authors",
		"ISBN",
		"ISBN13",
		"My Rating",
		"Average Rating",
		"Publisher",
		"Binding",
		"Number of Pages",
		"Year Published",
		"Original Publication Year",
		"Date Read",
		"Date Added",
		"Bookshelves",
		"Bookshelves with positions",
		"Exclusive Shelf",
		"My Review",
		"Spoiler",
		"Private Notes",
		"Read Count",
		"Owned Copies",
	}
}

func (f *GoodreadsFormat) ExportBook(book *blef.Book, entry *blef.Entry) []string {
	row := make([]string, len(f.GetExportHeaders()))

	// Book Id - use a generated ID if UUID, otherwise ISBN13
	bookID := book.ID
	if len(bookID) > 13 {
		bookID = "" // Don't export UUIDs as Book Id
	}
	row[0] = bookID

	// Title
	row[1] = book.Title

	// Author
	if len(book.Authors) > 0 {
		row[2] = book.Authors[0].Name
		// Author l-f (Last, First)
		row[3] = formatAuthorLastFirst(book.Authors[0].Name)

		// Additional Authors
		if len(book.Authors) > 1 {
			additionalAuthors := make([]string, len(book.Authors)-1)
			for i, author := range book.Authors[1:] {
				additionalAuthors[i] = author.Name
			}
			row[4] = strings.Join(additionalAuthors, ", ")
		}
	}

	// ISBNs - wrap in Excel formula to preserve leading zeros
	if book.Identifiers.ISBN10 != "" {
		row[5] = fmt.Sprintf(`="%s"`, book.Identifiers.ISBN10)
	}
	if book.Identifiers.ISBN13 != "" {
		row[6] = fmt.Sprintf(`=""%s""`, book.Identifiers.ISBN13)
	}

	// My Rating
	if entry != nil && entry.UserData.Rating > 0 {
		row[7] = fmt.Sprintf("%.0f", entry.UserData.Rating)
	} else {
		row[7] = "0"
	}

	// Average Rating - leave empty
	row[8] = ""

	// Edition info
	if book.Edition != nil {
		row[9] = book.Edition.Publisher
		row[11] = fmt.Sprintf("%d", book.Edition.Pages)
		row[12] = book.Edition.PublishedDate
	}

	// Binding - leave empty
	row[10] = ""

	// Original Publication Year - leave empty
	row[13] = ""

	// Entry data
	if entry != nil {
		// Date Read
		if len(entry.UserData.ReadDates) > 0 {
			if parsed, err := time.Parse("2006-01-02", entry.UserData.ReadDates[0].Finished); err == nil {
				row[14] = parsed.Format("2006/01/02")
			}
		}

		// Date Added
		if entry.UserData.AddedAt != nil {
			row[15] = entry.UserData.AddedAt.Format("2006/01/02")
		}

		// Bookshelves (tags)
		if len(entry.UserData.Tags) > 0 {
			row[16] = strings.Join(entry.UserData.Tags, ", ")
		}

		// Exclusive Shelf (status)
		row[18] = mapStatusToGoodreads(entry.UserData.Status)

		// My Review
		row[19] = entry.UserData.Review
	}

	// Leave empty: Bookshelves with positions, Spoiler, Private Notes, Read Count, Owned Copies
	row[17] = ""
	row[20] = ""
	row[21] = ""
	row[22] = ""
	row[23] = ""

	return row
}

// formatAuthorLastFirst converts "First Last" to "Last, First"
func formatAuthorLastFirst(name string) string {
	parts := strings.Fields(name)
	if len(parts) < 2 {
		return name
	}
	// Simple heuristic: last word is last name
	lastName := parts[len(parts)-1]
	firstName := strings.Join(parts[:len(parts)-1], " ")
	return fmt.Sprintf("%s, %s", lastName, firstName)
}

// mapStatusToGoodreads converts BLEF status to Goodreads shelf
func mapStatusToGoodreads(status string) string {
	switch status {
	case "read":
		return "read"
	case "reading":
		return "currently-reading"
	case "to-read":
		return "to-read"
	default:
		return "to-read"
	}
}
