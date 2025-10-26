package csv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCSV(t *testing.T) {
	// Create a temporary CSV file
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "test.csv")

	content := `Title,Author,ISBN
The Little Prince,Antoine de Saint-Exupéry,9780156013987
1984,George Orwell,9780451524935
`

	if err := os.WriteFile(csvFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	data, err := ParseCSV(csvFile)
	if err != nil {
		t.Fatalf("ParseCSV failed: %v", err)
	}

	if len(data.Headers) != 3 {
		t.Errorf("Expected 3 headers, got %d", len(data.Headers))
	}

	if len(data.Rows) != 2 {
		t.Errorf("Expected 2 rows, got %d", len(data.Rows))
	}

	if data.Headers[0] != "Title" {
		t.Errorf("Expected first header to be 'Title', got '%s'", data.Headers[0])
	}
}

func TestGetColumnIndex(t *testing.T) {
	data := &CSVData{
		Headers: []string{"Title", "Author", "ISBN"},
	}

	tests := []struct {
		name     string
		expected int
	}{
		{"Title", 0},
		{"title", 0}, // Case insensitive
		{"Author", 1},
		{"ISBN", 2},
		{"NonExistent", -1},
	}

	for _, tt := range tests {
		result := data.GetColumnIndex(tt.name)
		if result != tt.expected {
			t.Errorf("GetColumnIndex(%s) = %d, want %d", tt.name, result, tt.expected)
		}
	}
}

func TestGetValue(t *testing.T) {
	data := &CSVData{
		Headers: []string{"Title", "Author", "ISBN"},
	}

	row := []string{"Test Book", "Test Author", "1234567890"}

	tests := []struct {
		column   string
		expected string
	}{
		{"Title", "Test Book"},
		{"Author", "Test Author"},
		{"ISBN", "1234567890"},
		{"NonExistent", ""},
	}

	for _, tt := range tests {
		result := data.GetValue(row, tt.column)
		if result != tt.expected {
			t.Errorf("GetValue(row, %s) = %s, want %s", tt.column, result, tt.expected)
		}
	}
}
