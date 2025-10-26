package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set at build time via ldflags
	Version = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "blef-cli",
	Short: "BLEF - Book Library Exchange Format CLI tool",
	Long: `A command-line tool for working with BLEF (Book Library Exchange Format) files.

BLEF is an open, interoperable standard for exchanging personal book library data
between reading platforms, management tools, and online services.

Commands:
  validate - Validate a BLEF file against the JSON schema
  convert  - Convert CSV files to BLEF format
  view     - Interactive viewer for BLEF files`,
	Version: Version,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`{{.Version}}
`)
}
