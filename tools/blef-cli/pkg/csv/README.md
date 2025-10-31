# CSV Import/Export Package

This package provides a flexible architecture for importing and exporting BLEF documents from/to various CSV formats.

## ✨ Features

- **Bidirectional**: Full import AND export support
- **Interface-based**: Easy to extend with new formats
- **Auto-detection**: Automatically recognizes CSV formats
- **Type-safe**: Compile-time guarantees
- **Well-tested**: Comprehensive test coverage

## Architecture

The package uses an **interface-based design** that makes it easy to add support for new CSV formats:

```
CSVFormat (interface)
    ├── GoodreadsFormat
    ├── BabelioFormat
    └── YourCustomFormat (easy to add!)
```

## Adding a New Format

To add support for a new CSV format, simply create a new file and implement the `CSVFormat` interface:

### Step 1: Create your format file

Create a new file like `pkg/csv/myformat_format.go`:

```go
package csv

import (
    "strings"
    "github.com/yoanbernabeu/BLEF/tools/blef-cli/pkg/blef"
)

// MyFormat implements CSVFormat for My Platform library exports
type MyFormat struct{}

// Name returns the format identifier
func (f *MyFormat) Name() string {
    return "myformat"
}

// Description returns a human-readable description
func (f *MyFormat) Description() string {
    return "My Platform library export"
}

// Detect checks if the CSV matches this format
func (f *MyFormat) Detect(data *CSVData) bool {
    // Check for format-specific columns
    requiredColumns := []string{"BookID", "BookTitle", "BookAuthor"}
    for _, col := range requiredColumns {
        if data.GetColumnIndex(col) < 0 {
            return false
        }
    }
    return true
}

// GetImportMapping returns the column mapping for import
func (f *MyFormat) GetImportMapping() ColumnMapping {
    return ColumnMapping{
        ISBN13:        "ISBN",
        Title:         "BookTitle",
        Author:        "BookAuthor",
        Publisher:     "Publisher",
        PublishedDate: "PublishYear",
        Pages:         "PageCount",
        Rating:        "MyRating",
        Review:        "MyNotes",
        Status:        "ReadingStatus",
        DateRead:      "FinishedDate",
        DateAdded:     "AddedDate",
        Shelf:         "ShelfName",
    }
}

// CleanValue cleans format-specific values
func (f *MyFormat) CleanValue(value string) string {
    // Remove platform-specific formatting if needed
    return strings.TrimSpace(value)
}

// MapStatus converts platform-specific status to BLEF status
func (f *MyFormat) MapStatus(value string) string {
    value = strings.TrimSpace(strings.ToLower(value))
    
    switch value {
    case "finished":
        return "read"
    case "reading-now":
        return "reading"
    case "want-to-read":
        return "to-read"
    case "did-not-finish":
        return "abandoned"
    default:
        return "to-read"
    }
}

// MapRating converts platform-specific rating to BLEF rating (0-5)
func (f *MyFormat) MapRating(value string) float64 {
    // Parse and normalize to 0-5 scale
    var rating float64
    if _, err := fmt.Sscanf(value, "%f", &rating); err == nil {
        // If platform uses 10-point scale, convert to 5-point
        if rating > 5 {
            rating = rating / 2
        }
        return rating
    }
    return 0
}

// GetExportHeaders returns CSV headers for export
func (f *MyFormat) GetExportHeaders() []string {
    return []string{
        "BookID",
        "BookTitle",
        "BookAuthor",
        "ISBN",
        "Publisher",
        "PublishYear",
        "PageCount",
        "MyRating",
        "MyNotes",
        "ReadingStatus",
        "FinishedDate",
        "AddedDate",
        "ShelfName",
    }
}

// ExportBook converts a BLEF book and entry to a CSV row
func (f *MyFormat) ExportBook(book *blef.Book, entry *blef.Entry) []string {
    row := make([]string, len(f.GetExportHeaders()))
    
    // Map BLEF fields to CSV columns
    row[0] = book.ID
    row[1] = book.Title
    if len(book.Authors) > 0 {
        row[2] = book.Authors[0].Name
    }
    row[3] = book.Identifiers.ISBN13
    
    if book.Edition != nil {
        row[4] = book.Edition.Publisher
        row[5] = book.Edition.PublishedDate
        row[6] = fmt.Sprintf("%d", book.Edition.Pages)
    }
    
    if entry != nil {
        row[7] = fmt.Sprintf("%.1f", entry.UserData.Rating)
        row[8] = entry.UserData.Review
        row[9] = mapStatusToMyFormat(entry.UserData.Status)
        
        if len(entry.UserData.ReadDates) > 0 {
            row[10] = entry.UserData.ReadDates[0].Finished
        }
        if entry.UserData.AddedAt != nil {
            row[11] = entry.UserData.AddedAt.Format("2006-01-02")
        }
        if len(entry.CollectionIDs) > 0 {
            row[12] = entry.CollectionIDs[0]
        }
    }
    
    return row
}

// Helper function for export
func mapStatusToMyFormat(status string) string {
    switch status {
    case "read":
        return "finished"
    case "reading":
        return "reading-now"
    case "to-read":
        return "want-to-read"
    case "abandoned":
        return "did-not-finish"
    default:
        return "want-to-read"
    }
}
```

### Step 2: Register your format

Add your format to the default registry in `pkg/csv/format.go`:

```go
func init() {
    DefaultRegistry = NewFormatRegistry()
    
    // Register built-in formats
    DefaultRegistry.Register(&GoodreadsFormat{})
    DefaultRegistry.Register(&BabelioFormat{})
    DefaultRegistry.Register(&MyFormat{})  // Add your format here!
}
```

### Step 3: Done! 🎉

Your format is now available:

```bash
# Auto-detection
blef-cli convert mybooks.csv

# Explicit format
blef-cli convert mybooks.csv -f myformat

# Export (when implemented)
blef-cli export library.blef.json -f myformat -o mybooks.csv
```

## Testing

Add tests for your format in a `*_test.go` file:

```go
func TestMyFormatDetect(t *testing.T) {
    data := &CSVData{
        Headers: []string{"BookID", "BookTitle", "BookAuthor"},
    }
    
    format := &MyFormat{}
    if !format.Detect(data) {
        t.Error("Should detect MyFormat CSV")
    }
}

func TestMyFormatMapStatus(t *testing.T) {
    format := &MyFormat{}
    tests := []struct {
        input    string
        expected string
    }{
        {"finished", "read"},
        {"reading-now", "reading"},
        {"want-to-read", "to-read"},
    }
    
    for _, tt := range tests {
        result := format.MapStatus(tt.input)
        if result != tt.expected {
            t.Errorf("MapStatus(%s) = %s, want %s", tt.input, result, tt.expected)
        }
    }
}
```

## Built-in Formats

### Goodreads
- File: `goodreads_format.go`
- Handles Excel formula formatting (`=""value""`)
- Maps Goodreads shelves to BLEF status

### Babelio
- File: `babelio_format.go`
- Supports French status names
- Standard CSV format (no special formatting)

## Interface Reference

### CSVFormat Interface

```go
type CSVFormat interface {
    Name() string
    Description() string
    Detect(data *CSVData) bool
    GetImportMapping() ColumnMapping
    CleanValue(value string) string
    MapStatus(value string) string
    MapRating(value string) float64
    GetExportHeaders() []string
    ExportBook(book *blef.Book, entry *blef.Entry) []string
}
```

### Key Methods

- **Name()**: Unique identifier for the format (lowercase, no spaces)
- **Description()**: Human-readable name shown to users
- **Detect()**: Auto-detection logic based on CSV columns
- **GetImportMapping()**: Maps CSV columns to BLEF fields
- **CleanValue()**: Removes platform-specific formatting
- **MapStatus()**: Converts platform status to BLEF status
- **MapRating()**: Normalizes ratings to 0-5 scale
- **GetExportHeaders()**: CSV column headers for export
- **ExportBook()**: Converts BLEF data to CSV row

## Benefits of Interface-Based Design

✅ **Easy to extend**: Add new formats without modifying existing code  
✅ **Type-safe**: Compile-time checks ensure correct implementation  
✅ **Testable**: Each format can be tested independently  
✅ **Maintainable**: Clear separation of concerns  
✅ **Pluggable**: Formats can be added/removed easily  

## Contributing

When contributing a new format:

1. Implement all interface methods
2. Add comprehensive tests
3. Update this README with format details
4. Add example CSV files to documentation

