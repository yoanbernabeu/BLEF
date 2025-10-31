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
	exportFormat     string
	exportOutputFile string
)

var exportCmd = &cobra.Command{
	Use:   "export [blef-file]",
	Short: "Export BLEF file to CSV format",
	Long: `Export a BLEF file to CSV format compatible with various platforms.

Supported export formats:
  - goodreads: Goodreads CSV format (with Excel formulas)
  - babelio: Babelio CSV format (French)

The exported CSV can be imported back into the respective platform.

Examples:
  blef-cli export library.blef.json -f goodreads
  blef-cli export library.blef.json -f babelio -o export.csv
  blef-cli export library.blef.json -f goodreads -o goodreads_import.csv`,
	Args: cobra.ExactArgs(1),
	Run:  runExport,
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "", "Export format (goodreads, babelio) [required]")
	exportCmd.Flags().StringVarP(&exportOutputFile, "output", "o", "", "Output CSV file path (default: input-format.csv)")
	_ = exportCmd.MarkFlagRequired("format")
}

func runExport(cmd *cobra.Command, args []string) {
	inputFile := args[0]

	// Determine output file
	if exportOutputFile == "" {
		base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
		exportOutputFile = fmt.Sprintf("%s-%s.csv", base, exportFormat)
	}

	fmt.Printf("📤 Exporting BLEF to CSV format\n")
	fmt.Printf("Input:  %s\n", inputFile)
	fmt.Printf("Output: %s\n", exportOutputFile)
	fmt.Printf("Format: %s\n\n", exportFormat)

	// Load BLEF document
	fmt.Println("📖 Reading BLEF file...")
	doc, err := blef.LoadFromFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error reading BLEF file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✅ Loaded %d books, %d entries\n\n", len(doc.Books), len(doc.Entries))

	// Get export format
	format := csv.DefaultRegistry.GetByName(strings.ToLower(exportFormat))
	if format == nil {
		fmt.Fprintf(os.Stderr, "❌ Unknown export format: %s\n", exportFormat)
		fmt.Fprintf(os.Stderr, "Available formats: ")
		for i, f := range csv.DefaultRegistry.GetAll() {
			if i > 0 {
				fmt.Fprintf(os.Stderr, ", ")
			}
			fmt.Fprintf(os.Stderr, "%s", f.Name())
		}
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	// Create exporter
	exporter := csv.NewExporter(doc, format)

	// Show export stats
	stats := exporter.GetExportStats()
	fmt.Println("📊 Export preview:")
	fmt.Printf("  Total books:   %d\n", stats.TotalBooks)
	fmt.Printf("  Total entries: %d\n", stats.TotalEntries)
	fmt.Printf("  Will export:   %d rows\n", stats.Exported)
	if stats.Skipped > 0 {
		fmt.Printf("  ⚠️  Skipped:    %d entries (missing book data)\n", stats.Skipped)
	}
	fmt.Println()

	// Export to file
	fmt.Printf("💾 Writing to %s...\n", exportOutputFile)
	if err := exporter.ExportToFile(exportOutputFile); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Export failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Export complete!")
	fmt.Printf("\nYour CSV file is ready to import into %s.\n", format.Description())
}

