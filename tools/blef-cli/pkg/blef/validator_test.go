package blef

import (
	"testing"
	"time"
)

func TestValidateISBN13(t *testing.T) {
	tests := []struct {
		isbn  string
		valid bool
	}{
		{"9780156013987", true},
		{"9780062316097", true},
		{"9780439139595", true},
		{"1234567890123", false},
		{"978015601398X", false},
		{"978", false},
	}

	for _, tt := range tests {
		result := validateISBN13(tt.isbn)
		if result != tt.valid {
			t.Errorf("validateISBN13(%s) = %v, want %v", tt.isbn, result, tt.valid)
		}
	}
}

func TestValidateDocument(t *testing.T) {
	// Valid document
	validDoc := &BLEFDocument{
		Format:     "BLEF",
		Version:    "0.1.0",
		ExportedAt: time.Now(),
		Books: []Book{
			{
				ID:    "9780156013987",
				Title: "Test Book",
				Authors: []Author{
					{Name: "Test Author"},
				},
				Identifiers: Identifiers{
					ISBN13: "9780156013987",
				},
			},
		},
		Collections: []Collection{
			{
				ID:       "read",
				Name:     "Read",
				Type:     "read",
				IsPublic: true,
			},
		},
		Entries: []Entry{
			{
				BookID:        "9780156013987",
				CollectionIDs: []string{"read"},
				UserData: UserData{
					Status: "read",
					Rating: 4.5,
				},
			},
		},
	}

	errors := ValidateDocument(validDoc)
	if len(errors) > 0 {
		t.Errorf("ValidateDocument returned errors for valid document: %v", errors)
	}

	// Invalid format
	invalidFormat := *validDoc
	invalidFormat.Format = "WRONG"
	errors = ValidateDocument(&invalidFormat)
	if len(errors) == 0 {
		t.Error("ValidateDocument should return error for invalid format")
	}

	// Missing collections
	missingCollections := *validDoc
	missingCollections.Collections = []Collection{}
	errors = ValidateDocument(&missingCollections)
	if len(errors) == 0 {
		t.Error("ValidateDocument should return error for missing collections")
	}

	// Invalid ISBN
	invalidISBN := *validDoc
	invalidISBN.Books[0].ID = "1234567890123"
	errors = ValidateDocument(&invalidISBN)
	if len(errors) == 0 {
		t.Error("ValidateDocument should return error for invalid ISBN")
	}

	// Invalid rating range
	invalidRating := *validDoc
	invalidRating.Entries[0].UserData.Rating = 10
	errors = ValidateDocument(&invalidRating)
	if len(errors) == 0 {
		t.Error("ValidateDocument should return error for invalid rating")
	}
}

func TestCheckReferentialIntegrity(t *testing.T) {
	doc := &BLEFDocument{
		Books: []Book{
			{ID: "9780156013987", Title: "Book 1", Authors: []Author{{Name: "Author"}}, Identifiers: Identifiers{}},
		},
		Collections: []Collection{
			{ID: "read", Name: "Read", Type: "read"},
		},
		Entries: []Entry{
			{
				BookID:        "9780062316097", // Non-existent book
				CollectionIDs: []string{"read"},
				UserData:      UserData{Status: "read"},
			},
		},
	}

	errors := CheckReferentialIntegrity(doc)
	if len(errors) == 0 {
		t.Error("CheckReferentialIntegrity should return error for non-existent book reference")
	}

	// Non-existent collection
	doc.Entries[0].BookID = "9780156013987"
	doc.Entries[0].CollectionIDs = []string{"non-existent"}
	errors = CheckReferentialIntegrity(doc)
	if len(errors) == 0 {
		t.Error("CheckReferentialIntegrity should return error for non-existent collection reference")
	}

	// Invalid status
	doc.Entries[0].CollectionIDs = []string{"read"}
	doc.Entries[0].UserData.Status = "invalid-status"
	errors = CheckReferentialIntegrity(doc)
	if len(errors) == 0 {
		t.Error("CheckReferentialIntegrity should return error for invalid status")
	}
}
