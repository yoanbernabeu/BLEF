package csv

import (
	"fmt"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/google/uuid"
	"github.com/yoanbernabeu/BLEF/tools/blef-cli/pkg/blef"
)

// Mapper converts CSV data to BLEF format
type Mapper struct {
	Data    *CSVData
	Mapping ColumnMapping
	Format  CSVFormat
}

// NewMapper creates a new CSV to BLEF mapper
func NewMapper(data *CSVData, format CSVFormat) *Mapper {
	mapping := ColumnMapping{}
	if format != nil {
		mapping = format.GetImportMapping()
	}

	return &Mapper{
		Data:    data,
		Mapping: mapping,
		Format:  format,
	}
}

// InteractiveMapping prompts the user to map CSV columns
func (m *Mapper) InteractiveMapping() error {
	fmt.Println("\n📋 Column Mapping")
	fmt.Println("Map your CSV columns to BLEF fields...")
	fmt.Println()

	// Prepare field options
	fieldOptions := []string{
		"Book ID (ISBN-13 or unique ID)",
		"ISBN-13",
		"ISBN-10",
		"Title",
		"Author",
		"Language",
		"Publisher",
		"Published Date",
		"Number of Pages",
		"My Rating",
		"My Review",
		"Reading Status",
		"Date Read",
		"Date Added",
		"Tags",
		"Shelf/Collection",
		"(skip this column)",
	}

	// Map each CSV column
	for _, header := range m.Data.Headers {
		var selectedField string
		prompt := &survey.Select{
			Message: fmt.Sprintf("Map column '%s' to:", header),
			Options: fieldOptions,
			Default: m.guessField(header),
		}

		if err := survey.AskOne(prompt, &selectedField); err != nil {
			return fmt.Errorf("mapping cancelled: %w", err)
		}

		// Update mapping based on selection
		m.updateMapping(header, selectedField)
	}

	return nil
}

// guessField attempts to guess the appropriate BLEF field for a CSV column
func (m *Mapper) guessField(columnName string) string {
	lower := strings.ToLower(columnName)

	if strings.Contains(lower, "isbn") && strings.Contains(lower, "13") {
		return "ISBN-13"
	}
	if strings.Contains(lower, "isbn") {
		return "ISBN-10"
	}
	if strings.Contains(lower, "title") || strings.Contains(lower, "titre") {
		return "Title"
	}
	if strings.Contains(lower, "author") || strings.Contains(lower, "auteur") {
		return "Author"
	}
	if strings.Contains(lower, "rating") || strings.Contains(lower, "note") {
		return "My Rating"
	}
	if strings.Contains(lower, "review") || strings.Contains(lower, "critique") {
		return "My Review"
	}
	if strings.Contains(lower, "status") || strings.Contains(lower, "état") {
		return "Reading Status"
	}
	if strings.Contains(lower, "shelf") || strings.Contains(lower, "étagère") {
		return "Shelf/Collection"
	}
	if strings.Contains(lower, "language") || strings.Contains(lower, "langue") {
		return "Language"
	}
	if strings.Contains(lower, "publisher") || strings.Contains(lower, "editeur") {
		return "Publisher"
	}
	if strings.Contains(lower, "page") {
		return "Number of Pages"
	}
	if strings.Contains(lower, "tag") {
		return "Tags"
	}

	return "(skip this column)"
}

// updateMapping updates the mapping based on user selection
func (m *Mapper) updateMapping(columnName, selectedField string) {
	switch selectedField {
	case "Book ID (ISBN-13 or unique ID)":
		m.Mapping.BookID = columnName
	case "ISBN-13":
		m.Mapping.ISBN13 = columnName
	case "ISBN-10":
		m.Mapping.ISBN10 = columnName
	case "Title":
		m.Mapping.Title = columnName
	case "Author":
		m.Mapping.Author = columnName
	case "Language":
		m.Mapping.Language = columnName
	case "Publisher":
		m.Mapping.Publisher = columnName
	case "Published Date":
		m.Mapping.PublishedDate = columnName
	case "Number of Pages":
		m.Mapping.Pages = columnName
	case "My Rating":
		m.Mapping.Rating = columnName
	case "My Review":
		m.Mapping.Review = columnName
	case "Reading Status":
		m.Mapping.Status = columnName
	case "Date Read":
		m.Mapping.DateRead = columnName
	case "Date Added":
		m.Mapping.DateAdded = columnName
	case "Tags":
		m.Mapping.Tags = columnName
	case "Shelf/Collection":
		m.Mapping.Shelf = columnName
	}
}

// ConvertToBLEF converts CSV data to a BLEF document
func (m *Mapper) ConvertToBLEF() (*blef.BLEFDocument, error) {
	doc := blef.NewDocument()

	// Track collections
	collections := make(map[string]*blef.Collection)

	// Process each row
	for rowIdx, row := range m.Data.Rows {
		if len(row) == 0 {
			continue
		}

		// Build book
		book := m.buildBook(row, rowIdx)
		if book == nil {
			continue // Skip invalid rows
		}

		// Add book
		if err := doc.AddBook(*book); err != nil {
			// Book might already exist, that's ok
			if !strings.Contains(err.Error(), "already exists") {
				fmt.Printf("Warning: failed to add book at row %d: %v\n", rowIdx+2, err)
			}
		}

		// Build entry
		entry := m.buildEntry(row, book.ID, &collections)
		if entry != nil {
			// Ensure collections exist in document
			for collID, coll := range collections {
				existing := doc.GetCollectionByID(collID)
				if existing == nil {
					_ = doc.AddCollection(*coll)
				}
			}

			if err := doc.AddEntry(*entry); err != nil {
				fmt.Printf("Warning: failed to add entry at row %d: %v\n", rowIdx+2, err)
			}
		}
	}

	// Ensure at least one collection exists
	if len(doc.Collections) == 0 {
		_ = doc.AddCollection(blef.Collection{
			ID:       "default",
			Name:     "My Library",
			Type:     "custom",
			IsPublic: true,
		})
	}

	return doc, nil
}

// buildBook creates a book from a CSV row
func (m *Mapper) buildBook(row []string, rowIdx int) *blef.Book {
	title := m.getValue(row, m.Mapping.Title)
	if title == "" {
		return nil // Skip rows without title
	}

	// Get ISBN values and clean them using format-specific cleaning
	isbn13 := m.getValue(row, m.Mapping.ISBN13)
	isbn10 := m.getValue(row, m.Mapping.ISBN10)

	if m.Format != nil {
		isbn13 = m.Format.CleanValue(isbn13)
		isbn10 = m.Format.CleanValue(isbn10)
	}

	// Determine book ID - prioritize ISBN13, then generate UUID
	// Note: ISBN-10 is NOT valid as book ID in BLEF (only ISBN-13 or UUID)
	bookID := m.getValue(row, m.Mapping.BookID)
	if bookID == "" {
		bookID = isbn13
	}
	if bookID == "" {
		// Generate UUID if no ISBN-13 available (even if we have ISBN-10)
		bookID = uuid.New().String()
	}

	// Build identifiers
	identifiers := blef.Identifiers{}
	if isbn13 != "" {
		identifiers.ISBN13 = isbn13
	}
	if isbn10 != "" {
		identifiers.ISBN10 = isbn10
	}

	// Build author
	authorName := m.getValue(row, m.Mapping.Author)
	if authorName == "" {
		authorName = "Unknown"
	}

	authors := []blef.Author{
		{Name: authorName},
	}

	book := &blef.Book{
		ID:          bookID,
		Title:       title,
		Authors:     authors,
		Identifiers: identifiers,
	}

	// Optional fields
	if lang := m.getValue(row, m.Mapping.Language); lang != "" {
		book.Language = lang
	}

	// Edition info
	if publisher := m.getValue(row, m.Mapping.Publisher); publisher != "" ||
		m.getValue(row, m.Mapping.PublishedDate) != "" ||
		m.getValue(row, m.Mapping.Pages) != "" {

		edition := &blef.Edition{
			Publisher:     m.getValue(row, m.Mapping.Publisher),
			PublishedDate: m.getValue(row, m.Mapping.PublishedDate),
		}

		if pagesStr := m.getValue(row, m.Mapping.Pages); pagesStr != "" {
			var pages int
			if _, err := fmt.Sscanf(pagesStr, "%d", &pages); err == nil {
				edition.Pages = pages
			}
		}

		book.Edition = edition
	}

	return book
}

// buildEntry creates an entry from a CSV row
func (m *Mapper) buildEntry(row []string, bookID string, collections *map[string]*blef.Collection) *blef.Entry {
	// Determine status
	statusStr := m.getValue(row, m.Mapping.Status)
	status := "to-read" // default
	if m.Format != nil {
		status = m.Format.MapStatus(statusStr)
	} else {
		status = normalizeStatus(statusStr)
	}

	// Determine collection
	shelf := m.getValue(row, m.Mapping.Shelf)
	if shelf == "" {
		shelf = "default"
	}

	collectionID := strings.ToLower(strings.ReplaceAll(shelf, " ", "-"))
	if _, exists := (*collections)[collectionID]; !exists {
		collType := "custom"
		if strings.Contains(strings.ToLower(shelf), "read") {
			collType = "read"
		} else if strings.Contains(strings.ToLower(shelf), "reading") {
			collType = "reading"
		} else if strings.Contains(strings.ToLower(shelf), "to-read") {
			collType = "to-read"
		}

		(*collections)[collectionID] = &blef.Collection{
			ID:       collectionID,
			Name:     shelf,
			Type:     collType,
			IsPublic: true,
		}
	}

	// Build user data
	userData := blef.UserData{
		Status: status,
	}

	if ratingStr := m.getValue(row, m.Mapping.Rating); ratingStr != "" {
		if m.Format != nil {
			userData.Rating = m.Format.MapRating(ratingStr)
		} else {
			userData.Rating = parseRating(ratingStr)
		}
	}

	if review := m.getValue(row, m.Mapping.Review); review != "" {
		userData.Review = review
	}

	if tags := m.getValue(row, m.Mapping.Tags); tags != "" {
		userData.Tags = strings.Split(tags, ",")
	}

	if dateAdded := m.getValue(row, m.Mapping.DateAdded); dateAdded != "" {
		if t, err := parseDate(dateAdded); err == nil {
			userData.AddedAt = &t
		}
	}

	if dateRead := m.getValue(row, m.Mapping.DateRead); dateRead != "" {
		if t, err := parseDate(dateRead); err == nil {
			userData.ReadDates = []blef.ReadDate{
				{Finished: t.Format("2006-01-02")},
			}
		}
	}

	return &blef.Entry{
		BookID:        bookID,
		CollectionIDs: []string{collectionID},
		UserData:      userData,
	}
}

// getValue retrieves a value from the row using the mapping
func (m *Mapper) getValue(row []string, columnName string) string {
	if columnName == "" {
		return ""
	}
	return m.Data.GetValue(row, columnName)
}

// parseDate attempts to parse various date formats
func parseDate(dateStr string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"01/02/2006",
		"02/01/2006",
		time.RFC3339,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
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

// parseRating converts rating strings to float64
func parseRating(value string) float64 {
	value = strings.TrimSpace(value)
	if value == "" || value == "0" {
		return 0
	}

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
