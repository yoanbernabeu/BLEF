package blef

import (
	"embed"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed schema.json
var schemaFS embed.FS

var (
	isbn13Regex = regexp.MustCompile(`^97[89]\d{10}$`)
	uuidV4Regex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

// ValidateDocument performs comprehensive validation on a BLEF document
func ValidateDocument(doc *BLEFDocument) []error {
	var errors []error

	// Basic field validation
	if doc.Format != "BLEF" {
		errors = append(errors, ValidationError{Field: "format", Message: "must be 'BLEF'"})
	}

	if doc.Version == "" {
		errors = append(errors, ValidationError{Field: "version", Message: "is required"})
	}

	if len(doc.Collections) == 0 {
		errors = append(errors, ValidationError{Field: "collections", Message: "must contain at least one collection"})
	}

	// Validate book IDs
	bookIDs := make(map[string]bool)
	for i, book := range doc.Books {
		// Check for duplicate IDs
		if bookIDs[book.ID] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("books[%d].id", i),
				Message: fmt.Sprintf("duplicate book ID: %s", book.ID),
			})
		}
		bookIDs[book.ID] = true

		// Validate ID format
		if !isbn13Regex.MatchString(book.ID) && !uuidV4Regex.MatchString(book.ID) {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("books[%d].id", i),
				Message: "must be valid ISBN-13 or UUID v4",
			})
		}

		// Validate ISBN-13 check digit if applicable
		if isbn13Regex.MatchString(book.ID) {
			if !validateISBN13(book.ID) {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("books[%d].id", i),
					Message: "invalid ISBN-13 check digit",
				})
			}
		}

		// Validate required fields
		if book.Title == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("books[%d].title", i),
				Message: "is required",
			})
		}

		if len(book.Authors) == 0 {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("books[%d].authors", i),
				Message: "must contain at least one author",
			})
		}
	}

	// Validate collection IDs
	collectionIDs := make(map[string]bool)
	for i, collection := range doc.Collections {
		if collectionIDs[collection.ID] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("collections[%d].id", i),
				Message: fmt.Sprintf("duplicate collection ID: %s", collection.ID),
			})
		}
		collectionIDs[collection.ID] = true

		if collection.Name == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("collections[%d].name", i),
				Message: "is required",
			})
		}

		if collection.Type == "" {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("collections[%d].type", i),
				Message: "is required",
			})
		}
	}

	// Validate referential integrity
	refErrors := CheckReferentialIntegrity(doc)
	errors = append(errors, refErrors...)

	return errors
}

// CheckReferentialIntegrity validates that all references are valid
func CheckReferentialIntegrity(doc *BLEFDocument) []error {
	var errors []error

	bookIDs := make(map[string]bool)
	for _, book := range doc.Books {
		bookIDs[book.ID] = true
	}

	collectionIDs := make(map[string]bool)
	for _, collection := range doc.Collections {
		collectionIDs[collection.ID] = true
	}

	// Check entry references
	for i, entry := range doc.Entries {
		// Check book_id reference
		if !bookIDs[entry.BookID] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("entries[%d].book_id", i),
				Message: fmt.Sprintf("references non-existent book: %s", entry.BookID),
			})
		}

		// Check collection_ids references
		if len(entry.CollectionIDs) == 0 {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("entries[%d].collection_ids", i),
				Message: "must contain at least one collection",
			})
		}

		for j, collID := range entry.CollectionIDs {
			if !collectionIDs[collID] {
				errors = append(errors, ValidationError{
					Field:   fmt.Sprintf("entries[%d].collection_ids[%d]", i, j),
					Message: fmt.Sprintf("references non-existent collection: %s", collID),
				})
			}
		}

		// Validate status
		validStatuses := map[string]bool{
			"read": true, "reading": true, "to-read": true,
			"abandoned": true, "wishlist": true,
		}
		if !validStatuses[entry.UserData.Status] {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("entries[%d].user_data.status", i),
				Message: fmt.Sprintf("invalid status: %s", entry.UserData.Status),
			})
		}

		// Validate rating range
		if entry.UserData.Rating < 0 || entry.UserData.Rating > 5 {
			errors = append(errors, ValidationError{
				Field:   fmt.Sprintf("entries[%d].user_data.rating", i),
				Message: "must be between 0 and 5",
			})
		}
	}

	return errors
}

// ValidateAgainstSchema validates JSON data against the embedded JSON Schema
func ValidateAgainstSchema(jsonData []byte) error {
	schemaData, err := schemaFS.ReadFile("schema.json")
	if err != nil {
		return fmt.Errorf("failed to read embedded schema: %w", err)
	}

	schemaLoader := gojsonschema.NewBytesLoader(schemaData)
	documentLoader := gojsonschema.NewBytesLoader(jsonData)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if !result.Valid() {
		var errorMessages []string
		for _, err := range result.Errors() {
			errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", err.Field(), err.Description()))
		}
		return fmt.Errorf("schema validation errors:\n%s", strings.Join(errorMessages, "\n"))
	}

	return nil
}

// validateISBN13 validates the ISBN-13 check digit
func validateISBN13(isbn string) bool {
	if len(isbn) != 13 {
		return false
	}

	sum := 0
	for i, char := range isbn[:12] {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return false
		}
		if i%2 == 0 {
			sum += digit
		} else {
			sum += digit * 3
		}
	}

	checkDigit := (10 - (sum % 10)) % 10
	lastDigit, _ := strconv.Atoi(string(isbn[12]))

	return checkDigit == lastDigit
}
