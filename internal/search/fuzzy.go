package search

import (
	"sort"
	"strings"

	"github.com/nam-ruto/imperium/internal/models"
	"github.com/sahilm/fuzzy"
)

// commandSource wraps commands for fuzzy matching
type commandSource struct {
	commands []models.CommandRecord
}

func (c commandSource) String(i int) string {
	cmd := c.commands[i]
	// Combine tags, usage, and command for searching
	// Tags have highest priority by appearing first
	parts := []string{
		strings.Join(cmd.Tags, " "),
		cmd.Usage,
		cmd.Command,
	}
	if cmd.Notes != "" {
		parts = append(parts, cmd.Notes)
	}
	return strings.Join(parts, " ")
}

func (c commandSource) Len() int {
	return len(c.commands)
}

// Search performs fuzzy search on commands
func Search(query string, commands []models.CommandRecord) []models.CommandRecord {
	if query == "" {
		// Return all commands, sorted by favorites first, then recent
		return sortByRelevance(commands)
	}

	source := commandSource{commands: commands}
	matches := fuzzy.FindFrom(query, source)

	results := make([]models.CommandRecord, len(matches))
	for i, match := range matches {
		results[i] = commands[match.Index]
	}

	return results
}

// sortByRelevance sorts commands by favorites first, then by last used
func sortByRelevance(commands []models.CommandRecord) []models.CommandRecord {
	sorted := make([]models.CommandRecord, len(commands))
	copy(sorted, commands)

	sort.SliceStable(sorted, func(i, j int) bool {
		// Favorites first
		if sorted[i].Favorite != sorted[j].Favorite {
			return sorted[i].Favorite
		}
		// Then by last used (most recent first)
		if sorted[i].LastUsed != nil && sorted[j].LastUsed != nil {
			return sorted[i].LastUsed.After(*sorted[j].LastUsed)
		}
		if sorted[i].LastUsed != nil {
			return true
		}
		if sorted[j].LastUsed != nil {
			return false
		}
		// Finally by created_at (newest first)
		return sorted[i].CreatedAt.After(sorted[j].CreatedAt)
	})

	return sorted
}
