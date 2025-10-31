package csv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

// CSVData represents parsed CSV data
type CSVData struct {
	Headers []string
	Rows    [][]string
}

// ParseCSV reads and parses a CSV file with automatic encoding detection
func ParseCSV(filename string) (*CSVData, error) {
	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	// Auto-detect and convert encoding if needed
	content, err = convertToUTF8(content)
	if err != nil {
		return nil, fmt.Errorf("encoding conversion error: %w", err)
	}

	// Auto-detect delimiter
	delimiter := detectDelimiterFromContent(content)

	// Parse CSV from converted content
	reader := csv.NewReader(bytes.NewReader(content))
	reader.Comma = delimiter
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV headers: %w", err)
	}

	// Read all rows
	var rows [][]string
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read CSV row: %w", err)
		}
		rows = append(rows, row)
	}

	return &CSVData{
		Headers: headers,
		Rows:    rows,
	}, nil
}

// GetColumnIndex returns the index of a column by name (case-insensitive)
func (d *CSVData) GetColumnIndex(name string) int {
	lowerName := toLower(name)
	for i, header := range d.Headers {
		if toLower(header) == lowerName {
			return i
		}
	}
	return -1
}

// GetValue returns the value at the specified column for a given row
func (d *CSVData) GetValue(row []string, columnName string) string {
	idx := d.GetColumnIndex(columnName)
	if idx < 0 || idx >= len(row) {
		return ""
	}
	return row[idx]
}

func toLower(s string) string {
	return strings.ToLower(s)
}

// detectDelimiterFromContent auto-detects the CSV delimiter from content
func detectDelimiterFromContent(content []byte) rune {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	if !scanner.Scan() {
		return ',' // Default to comma
	}

	firstLine := scanner.Text()

	// Count occurrences of common delimiters
	delimiters := []rune{',', ';', '\t', '|'}
	maxCount := 0
	bestDelimiter := ','

	for _, delim := range delimiters {
		count := strings.Count(firstLine, string(delim))
		if count > maxCount {
			maxCount = count
			bestDelimiter = delim
		}
	}

	return bestDelimiter
}

// convertToUTF8 detects encoding and converts to UTF-8 if needed
func convertToUTF8(content []byte) ([]byte, error) {
	// Check if already valid UTF-8
	if utf8.Valid(content) {
		return content, nil
	}

	// Try Latin-1 (ISO-8859-1) - common for French exports
	decoder := charmap.ISO8859_1.NewDecoder()
	utf8Content, _, err := transform.Bytes(decoder, content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode from ISO-8859-1: %w", err)
	}

	return utf8Content, nil
}
