package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yoanbernabeu/BLEF/tools/blef-cli/pkg/blef"
	"github.com/yoanbernabeu/BLEF/tools/blef-cli/pkg/csv"
)

var (
	outputFile   string
	formatName   string
	skipValidate bool
)

var convertCmd = &cobra.Command{
	Use:   "convert [csv-file]",
	Short: "Convert CSV file to BLEF format",
	Long: `Convert a CSV file from various platforms to BLEF format.

Supported formats (auto-detected):
  - Goodreads library export
  - Babelio library export
  - Custom CSV (interactive mapping)

The tool will attempt to auto-detect the CSV format. If detection fails,
you will be prompted to manually map columns to BLEF fields.

Examples:
  blef-cli convert books.csv
  blef-cli convert books.csv -o my-library.blef.json
  blef-cli convert books.csv -f goodreads
  blef-cli convert books.csv --no-validate`,
	Args: cobra.ExactArgs(1),
	Run:  runConvert,
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path (default: input.blef.json)")
	convertCmd.Flags().StringVarP(&formatName, "format", "f", "", "Force format (goodreads, babelio)")
	convertCmd.Flags().BoolVar(&skipValidate, "no-validate", false, "Skip validation after conversion")
}

func runConvert(cmd *cobra.Command, args []string) {
	inputFile := args[0]

	// Determine output file
	if outputFile == "" {
		base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
		outputFile = base + ".blef.json"
	}

	fmt.Printf("📥 Converting CSV to BLEF format\n")
	fmt.Printf("Input:  %s\n", inputFile)
	fmt.Printf("Output: %s\n\n", outputFile)

	// Parse CSV
	fmt.Println("📖 Parsing CSV file...")
	data, err := csv.ParseCSV(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error parsing CSV: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Found %d rows with %d columns\n\n", len(data.Rows), len(data.Headers))

	// Detect or select preset
	var preset *csv.Preset
	if formatName != "" {
		// User specified format
		preset = findPreset(formatName)
		if preset == nil {
			fmt.Fprintf(os.Stderr, "❌ Unknown format: %s\n", formatName)
			os.Exit(1)
		}
		fmt.Printf("🎯 Using format: %s\n\n", preset.Description)
	} else {
		// Auto-detect
		fmt.Println("🔍 Detecting CSV format...")
		preset = csv.DetectPreset(data)
		if preset != nil {
			fmt.Printf("✅ Detected format: %s\n", preset.Description)
			fmt.Println("")
		} else {
			fmt.Println("⚠️  Could not auto-detect format")
			fmt.Println("")
		}
	}

	// Create mapper
	mapper := csv.NewMapper(data, preset)

	// If no preset or manual mapping requested, do interactive mapping
	if preset == nil {
		fmt.Println("📋 Starting interactive column mapping...")
		if err := mapper.InteractiveMapping(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Mapping error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("")
	}

	// Convert to BLEF
	fmt.Println("🔄 Converting to BLEF format...")
	doc, err := mapper.ConvertToBLEF()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Conversion error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Created BLEF document with %d books, %d collections, %d entries\n",
		len(doc.Books), len(doc.Collections), len(doc.Entries))
	fmt.Println("")

	// Validate before writing (unless skipped)
	if !skipValidate {
		fmt.Println("🔍 Validating BLEF document...")
		errors := blef.ValidateDocument(doc)
		if len(errors) > 0 {
			fmt.Println("⚠️  Validation warnings:")
			for _, err := range errors {
				fmt.Printf("  • %v\n", err)
			}
			fmt.Println("")
		} else {
			fmt.Println("✅ Validation passed")
			fmt.Println("")
		}
	}

	// Write to file
	fmt.Printf("💾 Writing to %s...\n", outputFile)
	jsonData, err := doc.ToJSON()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error generating JSON: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputFile, jsonData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Conversion complete!")
	fmt.Printf("\nYou can now validate your file with:\n  blef-cli validate %s\n", outputFile)
}

func findPreset(name string) *csv.Preset {
	name = strings.ToLower(name)
	for _, preset := range csv.AllPresets {
		if strings.ToLower(preset.Name) == name {
			return &preset
		}
	}
	return nil
}
