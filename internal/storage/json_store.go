package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/gamekeep/gamekeep/internal/models"
)

// MetadataStore defines the interface for storing and retrieving metadata
type MetadataStore interface {
	// Games
	SaveGames(games []models.Game) error
	LoadGames() ([]models.Game, error)
	
	// Checkpoints
	SaveCheckpoints(checkpoints []models.Checkpoint) error
	LoadCheckpoints() ([]models.Checkpoint, error)
}

// JSONStore implements MetadataStore using JSON files
type JSONStore struct {
	configDir         string
	gamesFile         string
	checkpointsFile   string
	mu                sync.RWMutex
}

// NewJSONStore creates a new JSON-based metadata store
func NewJSONStore(configDir string) (*JSONStore, error) {
	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	return &JSONStore{
		configDir:       configDir,
		gamesFile:       filepath.Join(configDir, "games.json"),
		checkpointsFile: filepath.Join(configDir, "checkpoints.json"),
	}, nil
}

// SaveGames saves games to JSON file with atomic write
func (s *JSONStore) SaveGames(games []models.Game) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.atomicWriteJSON(s.gamesFile, games)
}

// LoadGames loads games from JSON file
func (s *JSONStore) LoadGames() ([]models.Game, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var games []models.Game
	if err := s.readJSON(s.gamesFile, &games); err != nil {
		if os.IsNotExist(err) {
			return []models.Game{}, nil
		}
		return nil, err
	}
	return games, nil
}

// SaveCheckpoints saves checkpoints to JSON file with atomic write
func (s *JSONStore) SaveCheckpoints(checkpoints []models.Checkpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.atomicWriteJSON(s.checkpointsFile, checkpoints)
}

// LoadCheckpoints loads checkpoints from JSON file
func (s *JSONStore) LoadCheckpoints() ([]models.Checkpoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var checkpoints []models.Checkpoint
	if err := s.readJSON(s.checkpointsFile, &checkpoints); err != nil {
		if os.IsNotExist(err) {
			return []models.Checkpoint{}, nil
		}
		return nil, err
	}
	return checkpoints, nil
}

// atomicWriteJSON writes JSON data atomically using temp file + rename
func (s *JSONStore) atomicWriteJSON(filepath string, data interface{}) error {
	// Marshal with indentation for readability
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Create temporary file in the same directory
	tmpFile := filepath + ".tmp"
	
	// Write to temp file
	if err := os.WriteFile(tmpFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpFile, filepath); err != nil {
		os.Remove(tmpFile) // Clean up temp file on error
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}

// readJSON reads and unmarshals JSON from file
func (s *JSONStore) readJSON(filepath string, v interface{}) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}
