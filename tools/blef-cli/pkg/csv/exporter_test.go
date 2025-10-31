package csv

import (
	"os"
	"testing"

	"github.com/yoanbernabeu/BLEF/tools/blef-cli/pkg/blef"
)

func TestExporterStats(t *testing.T) {
	// Create a test document
	doc := blef.NewDocument()
	
	book1 := blef.Book{
		ID:    "9780123456789",
		Title: "Test Book 1",
		Authors: []blef.Author{
			{Name: "Test Author"},
		},
	}
	
	book2 := blef.Book{
		ID:    "9780987654321",
		Title: "Test Book 2",
		Authors: []blef.Author{
			{Name: "Another Author"},
		},
	}
	
	_ = doc.AddBook(book1)
	_ = doc.AddBook(book2)
	
	_ = doc.AddCollection(blef.Collection{
		ID:   "test",
		Name: "Test Collection",
		Type: "custom",
	})
	
	entry1 := blef.Entry{
		BookID:        "9780123456789",
		CollectionIDs: []string{"test"},
		UserData: blef.UserData{
			Status: "read",
			Rating: 5.0,
		},
	}
	
	entry2 := blef.Entry{
		BookID:        "9780987654321",
		CollectionIDs: []string{"test"},
		UserData: blef.UserData{
			Status: "reading",
		},
	}
	
	_ = doc.AddEntry(entry1)
	_ = doc.AddEntry(entry2)
	
	// Create exporter
	format := &GoodreadsFormat{}
	exporter := NewExporter(doc, format)
	
	// Get stats
	stats := exporter.GetExportStats()
	
	if stats.TotalBooks != 2 {
		t.Errorf("Expected 2 books, got %d", stats.TotalBooks)
	}
	
	if stats.TotalEntries != 2 {
		t.Errorf("Expected 2 entries, got %d", stats.TotalEntries)
	}
	
	if stats.Exported != 2 {
		t.Errorf("Expected 2 exported, got %d", stats.Exported)
	}
	
	if stats.Skipped != 0 {
		t.Errorf("Expected 0 skipped, got %d", stats.Skipped)
	}
}

func TestExportToFile(t *testing.T) {
	// Create a test document
	doc := blef.NewDocument()
	
	book := blef.Book{
		ID:    "9780123456789",
		Title: "Test Book",
		Authors: []blef.Author{
			{Name: "Test Author"},
		},
		Identifiers: blef.Identifiers{
			ISBN13: "9780123456789",
		},
	}
	
	_ = doc.AddBook(book)
	
	_ = doc.AddCollection(blef.Collection{
		ID:   "test",
		Name: "Test Collection",
		Type: "custom",
	})
	
	entry := blef.Entry{
		BookID:        "9780123456789",
		CollectionIDs: []string{"test"},
		UserData: blef.UserData{
			Status: "read",
			Rating: 4.5,
			Review: "Great book!",
		},
	}
	
	_ = doc.AddEntry(entry)
	
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "blef-export-test-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()
	
	// Export
	format := &GoodreadsFormat{}
	exporter := NewExporter(doc, format)
	
	if err := exporter.ExportToFile(tmpFile.Name()); err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	
	// Check file exists and has content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}
	
	if len(content) == 0 {
		t.Error("Exported file is empty")
	}
	
	// Basic content checks
	contentStr := string(content)
	if !containsSubstring(contentStr, "Test Book") {
		t.Error("Exported file should contain book title")
	}
	
	if !containsSubstring(contentStr, "Test Author") {
		t.Error("Exported file should contain author name")
	}
}

func TestExportBabelioFormat(t *testing.T) {
	// Create a test document
	doc := blef.NewDocument()
	
	book := blef.Book{
		ID:    "9780123456789",
		Title: "Livre Test",
		Authors: []blef.Author{
			{Name: "Auteur Test"},
		},
		Identifiers: blef.Identifiers{
			ISBN13: "9780123456789",
		},
	}
	
	_ = doc.AddBook(book)
	
	_ = doc.AddCollection(blef.Collection{
		ID:   "lu",
		Name: "Livres lus",
		Type: "read",
	})
	
	entry := blef.Entry{
		BookID:        "9780123456789",
		CollectionIDs: []string{"lu"},
		UserData: blef.UserData{
			Status: "read",
			Rating: 4.0,
		},
	}
	
	_ = doc.AddEntry(entry)
	
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "blef-export-babelio-test-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()
	
	// Export with Babelio format
	format := &BabelioFormat{}
	exporter := NewExporter(doc, format)
	
	if err := exporter.ExportToFile(tmpFile.Name()); err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	
	// Check file exists and has content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read exported file: %v", err)
	}
	
	if len(content) == 0 {
		t.Error("Exported file is empty")
	}
	
	// Check for Babelio headers (real format)
	contentStr := string(content)
	if !containsSubstring(contentStr, "ISBN") {
		t.Error("Babelio export should contain 'ISBN' header")
	}
	
	if !containsSubstring(contentStr, "Titre") {
		t.Error("Babelio export should contain 'Titre' header")
	}
	
	if !containsSubstring(contentStr, "Statut") {
		t.Error("Babelio export should contain 'Statut' header")
	}
	
	if !containsSubstring(contentStr, "Livre Test") {
		t.Error("Exported file should contain book title")
	}
}

// Helper function
func containsSubstring(haystack, needle string) bool {
	return len(haystack) >= len(needle) && 
		(haystack == needle || 
		 len(haystack) > len(needle) && 
		 (haystack[:len(needle)] == needle || 
		  containsSubstring(haystack[1:], needle)))
}

