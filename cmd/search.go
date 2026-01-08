package cmd

import (
	"strings"

	"github.com/nam-ruto/imperium/internal/tui"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search commands with optional pre-filled query",
	Long: `Search your stored commands using fuzzy matching.

Examples:
  imp search           # Open search with empty query
  imp search git       # Open search pre-filled with "git"
  imp search "docker ps"  # Open search pre-filled with "docker ps"`,
	Run: func(cmd *cobra.Command, args []string) {
		query := ""
		if len(args) > 0 {
			query = strings.Join(args, " ")
		}
		tui.RunSearch(store, query)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
