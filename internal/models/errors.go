package models

import "errors"

var (
	// Game errors
	ErrEmptyGameName  = errors.New("game name cannot be empty")
	ErrEmptySavePath  = errors.New("save path cannot be empty")
	ErrGameNotFound   = errors.New("game not found")
	ErrGameExists     = errors.New("game already exists")

	// Checkpoint errors
	ErrEmptyGameID          = errors.New("game ID cannot be empty")
	ErrEmptyCheckpointName  = errors.New("checkpoint name cannot be empty")
	ErrCheckpointNotFound   = errors.New("checkpoint not found")
	
	// Storage errors
	ErrInvalidPath          = errors.New("invalid path")
	ErrHashMismatch         = errors.New("hash mismatch")
)
