package csv

import "github.com/yoanbernabeu/BLEF/tools/blef-cli/pkg/blef"

// CSVFormat defines the interface for CSV import/export formats
// Implement this interface to add support for new CSV formats (e.g., Goodreads, Babelio)
type CSVFormat interface {
	// Name returns the format identifier (e.g., "goodreads", "babelio")
	Name() string

	// Description returns a human-readable description
	Description() string

	// Detect checks if the CSV data matches this format
	Detect(data *CSVData) bool

	// GetImportMapping returns the column mapping for import
	GetImportMapping() ColumnMapping

	// CleanValue cleans format-specific values (e.g., Excel formulas)
	CleanValue(value string) string

	// MapStatus converts format-specific status to BLEF status
	MapStatus(value string) string

	// MapRating converts format-specific rating to BLEF rating (0-5)
	MapRating(value string) float64

	// GetExportHeaders returns the CSV headers for export
	GetExportHeaders() []string

	// ExportBook converts a BLEF book and entry to a CSV row
	ExportBook(book *blef.Book, entry *blef.Entry) []string
}

// FormatRegistry manages available CSV formats
type FormatRegistry struct {
	formats []CSVFormat
}

// NewFormatRegistry creates a new format registry
func NewFormatRegistry() *FormatRegistry {
	return &FormatRegistry{
		formats: make([]CSVFormat, 0),
	}
}

// Register adds a format to the registry
func (r *FormatRegistry) Register(format CSVFormat) {
	r.formats = append(r.formats, format)
}

// GetByName returns a format by name
func (r *FormatRegistry) GetByName(name string) CSVFormat {
	for _, format := range r.formats {
		if format.Name() == name {
			return format
		}
	}
	return nil
}

// DetectFormat attempts to automatically detect the format
func (r *FormatRegistry) DetectFormat(data *CSVData) CSVFormat {
	for _, format := range r.formats {
		if format.Detect(data) {
			return format
		}
	}
	return nil
}

// GetAll returns all registered formats
func (r *FormatRegistry) GetAll() []CSVFormat {
	return r.formats
}

// DefaultRegistry is the global registry with built-in formats
var DefaultRegistry *FormatRegistry

func init() {
	DefaultRegistry = NewFormatRegistry()

	// Register built-in formats
	DefaultRegistry.Register(&GoodreadsFormat{})
	DefaultRegistry.Register(&BabelioFormat{})
}
