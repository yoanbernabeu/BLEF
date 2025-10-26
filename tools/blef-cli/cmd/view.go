package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/yoanbernabeu/BLEF/tools/blef-cli/pkg/blef"
	"github.com/yoanbernabeu/BLEF/tools/blef-cli/pkg/viewer"
)

var viewCmd = &cobra.Command{
	Use:   "view [file]",
	Short: "Interactive viewer for BLEF files",
	Long: `Launch an interactive terminal viewer for BLEF files.

Features:
  - Browse books, collections, and statistics
  - View detailed book information
  - Color-coded reading status
  - Keyboard navigation

Controls:
  ↑/↓ or j/k  - Navigate up/down
  Tab         - Switch between views (Books/Collections/Stats)
  Enter       - View book details
  Esc         - Go back
  q           - Quit`,
	Args: cobra.ExactArgs(1),
	Run:  runView,
}

func init() {
	rootCmd.AddCommand(viewCmd)
}

func runView(cmd *cobra.Command, args []string) {
	filename := args[0]

	// Read and parse file
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	doc, err := blef.FromJSON(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing BLEF file: %v\n", err)
		os.Exit(1)
	}

	// Create and run TUI
	model := viewer.NewModel(doc)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running viewer: %v\n", err)
		os.Exit(1)
	}
}
