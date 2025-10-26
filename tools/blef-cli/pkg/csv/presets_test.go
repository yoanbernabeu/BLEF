package csv

import "testing"

func TestDetectPreset(t *testing.T) {
	// Goodreads format
	goodreadsData := &CSVData{
		Headers: []string{"Book Id", "Title", "Author", "ISBN13", "My Rating"},
	}

	preset := DetectPreset(goodreadsData)
	if preset == nil {
		t.Error("DetectPreset should detect Goodreads format")
	} else if preset.Name != "goodreads" {
		t.Errorf("Expected 'goodreads', got '%s'", preset.Name)
	}

	// Babelio format
	babelioData := &CSVData{
		Headers: []string{"EAN", "Titre", "Auteur", "Étagère"},
	}

	preset = DetectPreset(babelioData)
	if preset == nil {
		t.Error("DetectPreset should detect Babelio format")
	} else if preset.Name != "babelio" {
		t.Errorf("Expected 'babelio', got '%s'", preset.Name)
	}

	// Unknown format
	unknownData := &CSVData{
		Headers: []string{"Unknown", "Columns"},
	}

	preset = DetectPreset(unknownData)
	if preset != nil {
		t.Error("DetectPreset should return nil for unknown format")
	}
}

func TestMapStatus(t *testing.T) {
	tests := []struct {
		input    string
		preset   *Preset
		expected string
	}{
		{"read", &GoodreadsPreset, "read"},
		{"currently-reading", &GoodreadsPreset, "reading"},
		{"to-read", &GoodreadsPreset, "to-read"},
		{"lu", &BabelioPreset, "read"},
		{"en cours", &BabelioPreset, "reading"},
		{"à lire", &BabelioPreset, "to-read"},
		{"abandonné", &BabelioPreset, "abandoned"},
		{"read", nil, "read"},
		{"reading now", nil, "reading"},
	}

	for _, tt := range tests {
		result := MapStatus(tt.input, tt.preset)
		if result != tt.expected {
			presetName := "nil"
			if tt.preset != nil {
				presetName = tt.preset.Name
			}
			t.Errorf("MapStatus(%s, %s) = %s, want %s", tt.input, presetName, result, tt.expected)
		}
	}
}

func TestMapRating(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5", 5.0},
		{"4.5", 4.5},
		{"0", 0.0},
		{"", 0.0},
		{"10", 5.0}, // 10-point scale converted to 5
		{"8", 4.0},  // 8/2 = 4
		{"-1", 0.0}, // Negative clamped to 0
		{"invalid", 0.0},
	}

	for _, tt := range tests {
		result := MapRating(tt.input)
		if result != tt.expected {
			t.Errorf("MapRating(%s) = %.1f, want %.1f", tt.input, result, tt.expected)
		}
	}
}
