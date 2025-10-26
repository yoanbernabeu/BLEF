package blef

import (
	"encoding/json"
	"fmt"
	"time"
)

// NewDocument creates a new BLEF document with default values
func NewDocument() *BLEFDocument {
	return &BLEFDocument{
		Format:      "BLEF",
		Version:     "0.1.0",
		ExportedAt:  time.Now().UTC(),
		Books:       []Book{},
		Collections: []Collection{},
		Entries:     []Entry{},
	}
}

// AddBook adds a book to the document if it doesn't already exist
func (d *BLEFDocument) AddBook(book Book) error {
	// Check for duplicate IDs
	for _, b := range d.Books {
		if b.ID == book.ID {
			return fmt.Errorf("book with ID %s already exists", book.ID)
		}
	}
	d.Books = append(d.Books, book)
	return nil
}

// AddCollection adds a collection to the document
func (d *BLEFDocument) AddCollection(collection Collection) error {
	// Check for duplicate IDs
	for _, c := range d.Collections {
		if c.ID == collection.ID {
			return fmt.Errorf("collection with ID %s already exists", collection.ID)
		}
	}
	d.Collections = append(d.Collections, collection)
	return nil
}

// AddEntry adds an entry to the document
func (d *BLEFDocument) AddEntry(entry Entry) error {
	// Check if book exists
	bookExists := false
	for _, b := range d.Books {
		if b.ID == entry.BookID {
			bookExists = true
			break
		}
	}
	if !bookExists {
		return fmt.Errorf("book with ID %s does not exist", entry.BookID)
	}

	// Check if collections exist
	for _, collID := range entry.CollectionIDs {
		collExists := false
		for _, c := range d.Collections {
			if c.ID == collID {
				collExists = true
				break
			}
		}
		if !collExists {
			return fmt.Errorf("collection with ID %s does not exist", collID)
		}
	}

	d.Entries = append(d.Entries, entry)
	return nil
}

// ToJSON converts the document to JSON bytes with indentation
func (d *BLEFDocument) ToJSON() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}

// FromJSON parses a BLEF document from JSON bytes
func FromJSON(data []byte) (*BLEFDocument, error) {
	var doc BLEFDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse BLEF document: %w", err)
	}
	return &doc, nil
}

// GetBookByID retrieves a book by its ID
func (d *BLEFDocument) GetBookByID(id string) *Book {
	for i := range d.Books {
		if d.Books[i].ID == id {
			return &d.Books[i]
		}
	}
	return nil
}

// GetCollectionByID retrieves a collection by its ID
func (d *BLEFDocument) GetCollectionByID(id string) *Collection {
	for i := range d.Collections {
		if d.Collections[i].ID == id {
			return &d.Collections[i]
		}
	}
	return nil
}

// GetEntriesForBook retrieves all entries for a specific book
func (d *BLEFDocument) GetEntriesForBook(bookID string) []Entry {
	var entries []Entry
	for _, e := range d.Entries {
		if e.BookID == bookID {
			entries = append(entries, e)
		}
	}
	return entries
}
