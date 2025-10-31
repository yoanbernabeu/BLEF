package csv

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/yoanbernabeu/BLEF/tools/blef-cli/pkg/blef"
)

// Exporter converts BLEF documents to CSV format
type Exporter struct {
	Document *blef.BLEFDocument
	Format   CSVFormat
}

// NewExporter creates a new BLEF to CSV exporter
func NewExporter(doc *blef.BLEFDocument, format CSVFormat) *Exporter {
	return &Exporter{
		Document: doc,
		Format:   format,
	}
}

// ExportToFile exports the BLEF document to a CSV file
func (e *Exporter) ExportToFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	headers := e.Format.GetExportHeaders()
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	// Create a map of book IDs to books for quick lookup
	bookMap := make(map[string]*blef.Book)
	for i := range e.Document.Books {
		bookMap[e.Document.Books[i].ID] = &e.Document.Books[i]
	}

	// Export each entry with its book
	for i := range e.Document.Entries {
		entry := &e.Document.Entries[i]
		book, exists := bookMap[entry.BookID]
		if !exists {
			// Skip entries without corresponding books
			continue
		}

		row := e.Format.ExportBook(book, entry)
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	return writer.Error()
}

// ExportStats returns statistics about the export
type ExportStats struct {
	TotalBooks   int
	TotalEntries int
	Exported     int
	Skipped      int
}

// GetExportStats returns statistics about what will be exported
func (e *Exporter) GetExportStats() ExportStats {
	stats := ExportStats{
		TotalBooks:   len(e.Document.Books),
		TotalEntries: len(e.Document.Entries),
	}

	// Create a map of book IDs
	bookMap := make(map[string]bool)
	for i := range e.Document.Books {
		bookMap[e.Document.Books[i].ID] = true
	}

	// Count exportable entries
	for i := range e.Document.Entries {
		if bookMap[e.Document.Entries[i].BookID] {
			stats.Exported++
		} else {
			stats.Skipped++
		}
	}

	return stats
}

