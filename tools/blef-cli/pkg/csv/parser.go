package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// CSVData represents parsed CSV data
type CSVData struct {
	Headers []string
	Rows    [][]string
}

// ParseCSV reads and parses a CSV file
func ParseCSV(filename string) (*CSVData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

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
