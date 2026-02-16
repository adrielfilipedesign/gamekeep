package models

import (
	"time"
)

// Game represents a registered game in the system
type Game struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	SavePath string `json:"save_path"`
}

// Checkpoint represents a save state snapshot
type Checkpoint struct {
	ID        string    `json:"id"`
	GameID    string    `json:"game_id"`
	Name      string    `json:"name"`
	Note      string    `json:"note"`
	VaultFile string    `json:"vault_file"`
	Hash      string    `json:"hash"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate validates game fields
func (g *Game) Validate() error {
	if g.Name == "" {
		return ErrEmptyGameName
	}
	if g.SavePath == "" {
		return ErrEmptySavePath
	}
	return nil
}

// Validate validates checkpoint fields
func (c *Checkpoint) Validate() error {
	if c.GameID == "" {
		return ErrEmptyGameID
	}
	if c.Name == "" {
		return ErrEmptyCheckpointName
	}
	return nil
}
