package core

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/adrielfilipedesign/gamekeep/internal/models"
	"github.com/adrielfilipedesign/gamekeep/internal/storage"
	"github.com/adrielfilipedesign/gamekeep/internal/vault"
)

// Service handles core business logic
type Service struct {
	store       storage.MetadataStore
	vaultMgr    *vault.Manager
}

// NewService creates a new service instance
func NewService(store storage.MetadataStore, vaultMgr *vault.Manager) *Service {
	return &Service{
		store:    store,
		vaultMgr: vaultMgr,
	}
}

// AddGame registers a new game in the system
func (s *Service) AddGame(name, savePath string) (*models.Game, error) {
	// Validate input
	if name == "" {
		return nil, models.ErrEmptyGameName
	}
	if savePath == "" {
		return nil, models.ErrEmptySavePath
	}

	// Clean and validate path
	cleanPath := filepath.Clean(savePath)

	// Load existing games
	games, err := s.store.LoadGames()
	if err != nil {
		return nil, fmt.Errorf("failed to load games: %w", err)
	}

	// Check for duplicate name
	for _, g := range games {
		if strings.EqualFold(g.Name, name) {
			return nil, models.ErrGameExists
		}
	}

	// Create game with sanitized ID
	game := &models.Game{
		ID:       s.sanitizeID(name),
		Name:     name,
		SavePath: cleanPath,
	}

	if err := game.Validate(); err != nil {
		return nil, err
	}

	// Add to list
	games = append(games, *game)

	// Save
	if err := s.store.SaveGames(games); err != nil {
		return nil, fmt.Errorf("failed to save games: %w", err)
	}

	return game, nil
}

// GetGame retrieves a game by ID or name
func (s *Service) GetGame(identifier string) (*models.Game, error) {
	games, err := s.store.LoadGames()
	if err != nil {
		return nil, fmt.Errorf("failed to load games: %w", err)
	}

	// Try exact ID match first
	for _, g := range games {
		if g.ID == identifier {
			return &g, nil
		}
	}

	// Try case-insensitive name match
	for _, g := range games {
		if strings.EqualFold(g.Name, identifier) {
			return &g, nil
		}
	}

	// Try partial match
	for _, g := range games {
		if strings.Contains(strings.ToLower(g.ID), strings.ToLower(identifier)) {
			return &g, nil
		}
	}

	return nil, models.ErrGameNotFound
}

// ListGames returns all registered games
func (s *Service) ListGames() ([]models.Game, error) {
	return s.store.LoadGames()
}

// CreateCheckpoint creates a new checkpoint for a game
func (s *Service) CreateCheckpoint(gameIdentifier, name, note string) (*models.Checkpoint, error) {
	// Get game
	game, err := s.GetGame(gameIdentifier)
	if err != nil {
		return nil, err
	}

	// Generate checkpoint ID
	checkpointID := uuid.New().String()

	// Create vault archive
	vaultFile, hash, err := s.vaultMgr.CreateCheckpoint(game.ID, checkpointID, game.SavePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkpoint archive: %w", err)
	}

	// Create checkpoint metadata
	checkpoint := &models.Checkpoint{
		ID:        checkpointID,
		GameID:    game.ID,
		Name:      name,
		Note:      note,
		VaultFile: vaultFile,
		Hash:      hash,
		CreatedAt: time.Now().UTC(),
	}

	if err := checkpoint.Validate(); err != nil {
		// Clean up vault file on validation error
		s.vaultMgr.DeleteCheckpoint(vaultFile)
		return nil, err
	}

	// Load existing checkpoints
	checkpoints, err := s.store.LoadCheckpoints()
	if err != nil {
		// Clean up vault file on error
		s.vaultMgr.DeleteCheckpoint(vaultFile)
		return nil, fmt.Errorf("failed to load checkpoints: %w", err)
	}

	// Add new checkpoint
	checkpoints = append(checkpoints, *checkpoint)

	// Save
	if err := s.store.SaveCheckpoints(checkpoints); err != nil {
		// Clean up vault file on save error
		s.vaultMgr.DeleteCheckpoint(vaultFile)
		return nil, fmt.Errorf("failed to save checkpoints: %w", err)
	}

	return checkpoint, nil
}

// ListCheckpoints returns checkpoints for a specific game
func (s *Service) ListCheckpoints(gameIdentifier string) ([]models.Checkpoint, error) {
	// Get game to validate it exists
	game, err := s.GetGame(gameIdentifier)
	if err != nil {
		return nil, err
	}

	// Load all checkpoints
	allCheckpoints, err := s.store.LoadCheckpoints()
	if err != nil {
		return nil, fmt.Errorf("failed to load checkpoints: %w", err)
	}

	// Filter by game ID
	var checkpoints []models.Checkpoint
	for _, cp := range allCheckpoints {
		if cp.GameID == game.ID {
			checkpoints = append(checkpoints, cp)
		}
	}

	return checkpoints, nil
}

// GetCheckpoint retrieves a checkpoint by ID
func (s *Service) GetCheckpoint(checkpointID string) (*models.Checkpoint, error) {
	checkpoints, err := s.store.LoadCheckpoints()
	if err != nil {
		return nil, fmt.Errorf("failed to load checkpoints: %w", err)
	}

	for _, cp := range checkpoints {
		if cp.ID == checkpointID || strings.HasPrefix(cp.ID, checkpointID) {
			return &cp, nil
		}
	}

	return nil, models.ErrCheckpointNotFound
}

// RestoreCheckpoint restores a checkpoint to the game's save directory
func (s *Service) RestoreCheckpoint(checkpointID string) error {
	// Get checkpoint
	checkpoint, err := s.GetCheckpoint(checkpointID)
	if err != nil {
		return err
	}

	// Get game
	game, err := s.GetGame(checkpoint.GameID)
	if err != nil {
		return err
	}

	// Verify checkpoint integrity
	if err := s.vaultMgr.VerifyCheckpoint(checkpoint.VaultFile, checkpoint.Hash); err != nil {
		return fmt.Errorf("checkpoint verification failed: %w", err)
	}

	// Restore
	if err := s.vaultMgr.RestoreCheckpoint(checkpoint.VaultFile, game.SavePath); err != nil {
		return fmt.Errorf("failed to restore checkpoint: %w", err)
	}

	return nil
}

// DeleteCheckpoint removes a checkpoint
func (s *Service) DeleteCheckpoint(checkpointID string) error {
	// Get checkpoint
	checkpoint, err := s.GetCheckpoint(checkpointID)
	if err != nil {
		return err
	}

	// Load all checkpoints
	checkpoints, err := s.store.LoadCheckpoints()
	if err != nil {
		return fmt.Errorf("failed to load checkpoints: %w", err)
	}

	// Remove from list
	var updatedCheckpoints []models.Checkpoint
	for _, cp := range checkpoints {
		if cp.ID != checkpoint.ID {
			updatedCheckpoints = append(updatedCheckpoints, cp)
		}
	}

	// Save updated list
	if err := s.store.SaveCheckpoints(updatedCheckpoints); err != nil {
		return fmt.Errorf("failed to save checkpoints: %w", err)
	}

	// Delete vault file
	if err := s.vaultMgr.DeleteCheckpoint(checkpoint.VaultFile); err != nil {
		// Log but don't fail if vault file doesn't exist
		return fmt.Errorf("warning: failed to delete vault file: %w", err)
	}

	return nil
}

// sanitizeID creates a safe ID from a name
func (s *Service) sanitizeID(name string) string {
	// Convert to lowercase
	id := strings.ToLower(name)
	
	// Replace spaces and special chars with underscores
	id = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, id)
	
	// Remove consecutive underscores
	for strings.Contains(id, "__") {
		id = strings.ReplaceAll(id, "__", "_")
	}
	
	// Trim underscores from edges
	id = strings.Trim(id, "_")
	
	// Limit length
	if len(id) > 50 {
		id = id[:50]
	}
	
	return id
}
