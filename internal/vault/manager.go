package vault

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Manager handles vault operations for save files
type Manager struct {
	vaultDir string
}

// NewManager creates a new vault manager
func NewManager(vaultDir string) (*Manager, error) {
	// Ensure vault directory exists
	if err := os.MkdirAll(vaultDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create vault directory: %w", err)
	}

	return &Manager{
		vaultDir: vaultDir,
	}, nil
}

// CreateCheckpoint creates a compressed archive of the save directory
func (m *Manager) CreateCheckpoint(gameID, checkpointID, savePath string) (vaultFile string, hash string, err error) {
	// Create game-specific directory
	gameVaultDir := filepath.Join(m.vaultDir, gameID)
	if err := os.MkdirAll(gameVaultDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create game vault directory: %w", err)
	}

	// Destination zip file
	zipPath := filepath.Join(gameVaultDir, checkpointID+".zip")

	// Verify source path exists
	if _, err := os.Stat(savePath); err != nil {
		return "", "", fmt.Errorf("save path does not exist: %w", err)
	}

	// Create zip archive
	if err := m.zipDirectory(savePath, zipPath); err != nil {
		return "", "", fmt.Errorf("failed to create zip archive: %w", err)
	}

	// Calculate SHA256 hash
	hash, err = m.calculateHash(zipPath)
	if err != nil {
		os.Remove(zipPath) // Clean up on error
		return "", "", fmt.Errorf("failed to calculate hash: %w", err)
	}

	// Return relative path from vault root
	relPath, err := filepath.Rel(m.vaultDir, zipPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to get relative path: %w", err)
	}

	return relPath, hash, nil
}

// RestoreCheckpoint extracts a checkpoint archive to the save directory
func (m *Manager) RestoreCheckpoint(vaultFile, targetPath string) error {
	// Full path to zip file
	zipPath := filepath.Join(m.vaultDir, vaultFile)

	// Verify zip exists
	if _, err := os.Stat(zipPath); err != nil {
		return fmt.Errorf("checkpoint file not found: %w", err)
	}

	// Remove existing save directory
	if err := os.RemoveAll(targetPath); err != nil {
		return fmt.Errorf("failed to remove existing save directory: %w", err)
	}

	// Create target directory
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Extract zip
	if err := m.unzipArchive(zipPath, targetPath); err != nil {
		return fmt.Errorf("failed to extract checkpoint: %w", err)
	}

	return nil
}

// VerifyCheckpoint verifies the integrity of a checkpoint
func (m *Manager) VerifyCheckpoint(vaultFile, expectedHash string) error {
	zipPath := filepath.Join(m.vaultDir, vaultFile)

	actualHash, err := m.calculateHash(zipPath)
	if err != nil {
		return fmt.Errorf("failed to calculate hash: %w", err)
	}

	if actualHash != expectedHash {
		return fmt.Errorf("hash mismatch: expected %s, got %s", expectedHash, actualHash)
	}

	return nil
}

// DeleteCheckpoint removes a checkpoint file from the vault
func (m *Manager) DeleteCheckpoint(vaultFile string) error {
	zipPath := filepath.Join(m.vaultDir, vaultFile)
	return os.Remove(zipPath)
}

// zipDirectory compresses a directory into a zip file
func (m *Manager) zipDirectory(sourceDir, targetZip string) error {
	// Create zip file
	zipFile, err := os.Create(targetZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	// Walk through directory
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if path == sourceDir {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Create zip header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Use forward slashes for cross-platform compatibility
		header.Name = filepath.ToSlash(relPath)

		// Handle directories
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		// Create entry
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		// Write file content
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(writer, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// unzipArchive extracts a zip file to a target directory
func (m *Manager) unzipArchive(zipPath, targetDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		// Prevent zip slip vulnerability
		filePath := filepath.Join(targetDir, file.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(targetDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, file.Mode())
			continue
		}

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return err
		}

		// Extract file
		if err := m.extractFile(file, filePath); err != nil {
			return err
		}
	}

	return nil
}

// extractFile extracts a single file from zip archive
func (m *Manager) extractFile(file *zip.File, targetPath string) error {
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// calculateHash calculates SHA256 hash of a file
func (m *Manager) calculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
