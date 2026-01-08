package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/nam-ruto/imperium/internal/models"
)

// Storage defines the interface for command persistence
type Storage interface {
	Load() ([]models.CommandRecord, error)
	Save(commands []models.CommandRecord) error
	Add(command models.CommandRecord) error
}

// JSONStorage implements Storage using a JSON file
type JSONStorage struct {
	path string
}

// DefaultPath returns the default storage path
func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", "imperium", "commands.json")
}

// NewJSONStorage creates a new JSON storage instance
func NewJSONStorage(path string) *JSONStorage {
	if path == "" {
		path = DefaultPath()
	}
	return &JSONStorage{path: path}
}

// ensureDir creates the directory if it doesn't exist
func (s *JSONStorage) ensureDir() error {
	dir := filepath.Dir(s.path)
	return os.MkdirAll(dir, 0755)
}

// Load reads all commands from the JSON file
func (s *JSONStorage) Load() ([]models.CommandRecord, error) {
	if err := s.ensureDir(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.CommandRecord{}, nil
		}
		return nil, err
	}

	if len(data) == 0 {
		return []models.CommandRecord{}, nil
	}

	var storage models.CommandStorage
	if err := json.Unmarshal(data, &storage); err != nil {
		return nil, err
	}

	return storage.Commands, nil
}

// Save writes all commands to the JSON file atomically
func (s *JSONStorage) Save(commands []models.CommandRecord) error {
	if err := s.ensureDir(); err != nil {
		return err
	}

	storage := models.CommandStorage{Commands: commands}
	data, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return err
	}

	// Write to temp file first (atomic write)
	tmpFile := s.path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return err
	}

	// Atomic rename
	return os.Rename(tmpFile, s.path)
}

// Add appends a new command to storage
func (s *JSONStorage) Add(command models.CommandRecord) error {
	commands, err := s.Load()
	if err != nil {
		return err
	}

	commands = append(commands, command)
	return s.Save(commands)
}

// UpdateLastUsed updates the last_used field for a command
func (s *JSONStorage) UpdateLastUsed(commandText string) error {
	commands, err := s.Load()
	if err != nil {
		return err
	}

	for i := range commands {
		if commands[i].Command == commandText {
			commands[i].MarkUsed()
			break
		}
	}

	return s.Save(commands)
}
