package models

import "time"

// CommandRecord represents a stored command with metadata
type CommandRecord struct {
	Command  string   `json:"command"`
	Usage    string   `json:"usage"`
	Tags     []string `json:"tags"`
	Platform string   `json:"platform"` // macos | linux | windows

	// Optional fields
	Examples  []string   `json:"examples,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	Favorite  bool       `json:"favorite,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	LastUsed  *time.Time `json:"last_used,omitempty"`
}

// CommandStorage is the root structure for the JSON file
type CommandStorage struct {
	Commands []CommandRecord `json:"commands"`
}

// NewCommandRecord creates a new command record with defaults
func NewCommandRecord(command, usage string, tags []string, platform string) CommandRecord {
	if platform == "" {
		platform = "macos"
	}
	if tags == nil {
		tags = []string{}
	}
	return CommandRecord{
		Command:   command,
		Usage:     usage,
		Tags:      tags,
		Platform:  platform,
		CreatedAt: time.Now(),
	}
}

// MarkUsed updates the last_used timestamp
func (c *CommandRecord) MarkUsed() {
	now := time.Now()
	c.LastUsed = &now
}
