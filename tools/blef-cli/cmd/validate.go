package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yoanbernabeu/BLEF/tools/blef-cli/pkg/blef"
)

var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate a BLEF file against the JSON schema",
	Long: `Validate a BLEF file for correctness.

This command performs comprehensive validation including:
- JSON schema validation
- Referential integrity checks
- ISBN-13 check digit validation
- Required field validation
- Status and rating range validation

Exit codes:
  0 - File is valid
  1 - File is invalid or validation error`,
	Args: cobra.ExactArgs(1),
	Run:  runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) {
	filename := args[0]

	// Read file
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse BLEF document
	doc, err := blef.FromJSON(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error parsing BLEF file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📚 Validating BLEF file: %s\n\n", filename)

	// Validate against schema
	fmt.Println("🔍 Checking JSON schema...")
	if err := blef.ValidateAgainstSchema(data); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Schema validation failed:\n%v\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Schema validation passed")

	// Validate document structure and integrity
	fmt.Println("\n🔍 Checking document integrity...")
	errors := blef.ValidateDocument(doc)
	if len(errors) > 0 {
		fmt.Println("❌ Validation errors found:")
		for _, err := range errors {
			fmt.Printf("  • %v\n", err)
		}
		os.Exit(1)
	}
	fmt.Println("✅ Document integrity validated")

	// Display statistics
	fmt.Println("\n📊 Document Statistics:")
	fmt.Printf("  Format: %s v%s\n", doc.Format, doc.Version)
	fmt.Printf("  Exported: %s\n", doc.ExportedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("  Books: %d\n", len(doc.Books))
	fmt.Printf("  Collections: %d\n", len(doc.Collections))
	fmt.Printf("  Entries: %d\n", len(doc.Entries))

	if doc.User != nil && doc.User.Name != "" {
		fmt.Printf("  User: %s\n", doc.User.Name)
	}

	// Status breakdown
	statusCount := make(map[string]int)
	for _, entry := range doc.Entries {
		statusCount[entry.UserData.Status]++
	}

	if len(statusCount) > 0 {
		fmt.Println("\n📖 Reading Status:")
		for status, count := range statusCount {
			emoji := getStatusEmoji(status)
			fmt.Printf("  %s %s: %d\n", emoji, status, count)
		}
	}

	fmt.Println("\n✅ File is valid!")
}

func getStatusEmoji(status string) string {
	switch status {
	case "read":
		return "✅"
	case "reading":
		return "📖"
	case "to-read":
		return "📚"
	case "abandoned":
		return "❌"
	case "wishlist":
		return "⭐"
	default:
		return "📙"
	}
}
