package cmd

import (
	"os"

	"github.com/nam-ruto/imperium/internal/storage"
	"github.com/nam-ruto/imperium/internal/tui"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time
	Version = "dev"

	// Storage instance
	store *storage.JSONStorage
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "imp",
	Short: "Command your commands",
	Long: `Imperium - A fuzzy command-line snippet manager

Imperium helps you store, search, and recall terminal commands.
Think of it as a command palette for your terminal.

Usage:
  imp              Launch search mode (default)
  imp search git   Search with pre-filled query
  imp add          Add a new command
  imp browse       Browse all commands`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default action: launch search mode
		tui.RunSearch(store, "")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
	store = storage.NewJSONStorage("")
}
