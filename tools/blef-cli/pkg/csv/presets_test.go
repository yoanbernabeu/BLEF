package csv

import "testing"

func TestGoodreadsFormatDetect(t *testing.T) {
	goodreadsData := &CSVData{
		Headers: []string{"Book Id", "Title", "Author", "ISBN13", "My Rating"},
	}

	format := &GoodreadsFormat{}
	if !format.Detect(goodreadsData) {
		t.Error("GoodreadsFormat should detect Goodreads CSV")
	}

	invalidData := &CSVData{
		Headers: []string{"Unknown", "Columns"},
	}
	if format.Detect(invalidData) {
		t.Error("GoodreadsFormat should not detect invalid CSV")
	}
}

func TestBabelioFormatDetect(t *testing.T) {
	// Real Babelio export format
	babelioData := &CSVData{
		Headers: []string{"ISBN", "Titre", "Auteur", "Statut"},
	}

	format := &BabelioFormat{}
	if !format.Detect(babelioData) {
		t.Error("BabelioFormat should detect Babelio CSV")
	}

	invalidData := &CSVData{
		Headers: []string{"Unknown", "Columns"},
	}
	if format.Detect(invalidData) {
		t.Error("BabelioFormat should not detect invalid CSV")
	}
}

func TestFormatRegistry(t *testing.T) {
	registry := NewFormatRegistry()

	goodreads := &GoodreadsFormat{}
	babelio := &BabelioFormat{}

	registry.Register(goodreads)
	registry.Register(babelio)

	// Test GetByName
	found := registry.GetByName("goodreads")
	if found == nil {
		t.Error("Should find goodreads format")
	}
	if found.Name() != "goodreads" {
		t.Errorf("Expected 'goodreads', got '%s'", found.Name())
	}

	// Test DetectFormat
	goodreadsData := &CSVData{
		Headers: []string{"Book Id", "Title", "Author", "ISBN13", "My Rating"},
	}
	detected := registry.DetectFormat(goodreadsData)
	if detected == nil {
		t.Error("Should detect goodreads format")
	}
	if detected.Name() != "goodreads" {
		t.Errorf("Expected 'goodreads', got '%s'", detected.Name())
	}

	// Test GetAll
	all := registry.GetAll()
	if len(all) != 2 {
		t.Errorf("Expected 2 formats, got %d", len(all))
	}
}

func TestGoodreadsMapStatus(t *testing.T) {
	format := &GoodreadsFormat{}
	tests := []struct {
		input    string
		expected string
	}{
		{"read", "read"},
		{"currently-reading", "reading"},
		{"to-read", "to-read"},
		{"unknown", "to-read"},
	}

	for _, tt := range tests {
		result := format.MapStatus(tt.input)
		if result != tt.expected {
			t.Errorf("MapStatus(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestBabelioMapStatus(t *testing.T) {
	format := &BabelioFormat{}
	tests := []struct {
		input    string
		expected string
	}{
		{"lu", "read"},
		{"en cours", "reading"},
		{"à lire", "to-read"},
		{"abandonné", "abandoned"},
		{"unknown", "to-read"},
	}

	for _, tt := range tests {
		result := format.MapStatus(tt.input)
		if result != tt.expected {
			t.Errorf("MapStatus(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestGoodreadsMapRating(t *testing.T) {
	format := &GoodreadsFormat{}
	tests := []struct {
		input    string
		expected float64
	}{
		{"5", 5.0},
		{"4.5", 4.5},
		{"0", 0.0},
		{"", 0.0},
		{"-1", 0.0},
		{"invalid", 0.0},
	}

	for _, tt := range tests {
		result := format.MapRating(tt.input)
		if result != tt.expected {
			t.Errorf("MapRating(%s) = %.1f, want %.1f", tt.input, result, tt.expected)
		}
	}
}

func TestGoodreadsCleanValue(t *testing.T) {
	format := &GoodreadsFormat{}
	tests := []struct {
		input    string
		expected string
	}{
		{`=""9791028107819""`, "9791028107819"},
		{`="2367935947"`, "2367935947"},
		{`=""`, ""},
		{"9780123456789", "9780123456789"},
		{"", ""},
		{`  =""9780123456789""  `, "9780123456789"},
		{"normal value", "normal value"},
	}

	for _, tt := range tests {
		result := format.CleanValue(tt.input)
		if result != tt.expected {
			t.Errorf("CleanValue(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestDefaultRegistry(t *testing.T) {
	// Test that DefaultRegistry is initialized with built-in formats
	if DefaultRegistry == nil {
		t.Fatal("DefaultRegistry should not be nil")
	}

	formats := DefaultRegistry.GetAll()
	if len(formats) < 2 {
		t.Errorf("DefaultRegistry should have at least 2 formats, got %d", len(formats))
	}

	// Check that goodreads and babelio are registered
	goodreads := DefaultRegistry.GetByName("goodreads")
	if goodreads == nil {
		t.Error("DefaultRegistry should have goodreads format")
	}

	babelio := DefaultRegistry.GetByName("babelio")
	if babelio == nil {
		t.Error("DefaultRegistry should have babelio format")
	}
}
