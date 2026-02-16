package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/adrielfilipedesign/gamekeep/internal/core"
	"github.com/adrielfilipedesign/gamekeep/internal/storage"
	"github.com/adrielfilipedesign/gamekeep/internal/vault"
)

const (
	version = "1.0.0"
)

// CLI manages command-line interface
type CLI struct {
	service *core.Service
}

// NewCLI creates a new CLI instance
func NewCLI(service *core.Service) *CLI {
	return &CLI{
		service: service,
	}
}

// Run executes the CLI
func (c *CLI) Run(args []string) error {
	if len(args) < 1 {
		c.printUsage()
		return nil
	}

	command := args[0]
	
	switch command {
	case "add-game":
		return c.addGame(args[1:])
	case "list-games":
		return c.listGames()
	case "checkpoint":
		return c.createCheckpoint(args[1:])
	case "list":
		return c.listCheckpoints(args[1:])
	case "restore":
		return c.restoreCheckpoint(args[1:])
	case "delete":
		return c.deleteCheckpoint(args[1:])
	case "version":
		fmt.Printf("GameKeep v%s\n", version)
		return nil
	case "help":
		c.printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// addGame handles the add-game command
func (c *CLI) addGame(args []string) error {
	fs := flag.NewFlagSet("add-game", flag.ExitOnError)
	name := fs.String("name", "", "Game name (required)")
	path := fs.String("path", "", "Save directory path (required)")
	
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *name == "" || *path == "" {
		return fmt.Errorf("both --name and --path are required")
	}

	game, err := c.service.AddGame(*name, *path)
	if err != nil {
		return fmt.Errorf("failed to add game: %w", err)
	}

	fmt.Printf("✓ Game added successfully\n")
	fmt.Printf("  ID:   %s\n", game.ID)
	fmt.Printf("  Name: %s\n", game.Name)
	fmt.Printf("  Path: %s\n", game.SavePath)
	
	return nil
}

// listGames handles the list-games command
func (c *CLI) listGames() error {
	games, err := c.service.ListGames()
	if err != nil {
		return fmt.Errorf("failed to list games: %w", err)
	}

	if len(games) == 0 {
		fmt.Println("No games registered yet.")
		fmt.Println("Use 'gamekeep add-game --name \"Game Name\" --path \"/path/to/saves\"' to add one.")
		return nil
	}

	fmt.Printf("Registered Games (%d):\n\n", len(games))
	
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSAVE PATH")
	fmt.Fprintln(w, "──\t────\t─────────")
	
	for _, game := range games {
		fmt.Fprintf(w, "%s\t%s\t%s\n", game.ID, game.Name, game.SavePath)
	}
	
	w.Flush()
	return nil
}

// createCheckpoint handles the checkpoint command
func (c *CLI) createCheckpoint(args []string) error {
	fs := flag.NewFlagSet("checkpoint", flag.ExitOnError)
	game := fs.String("game", "", "Game ID or name (required)")
	name := fs.String("name", "", "Checkpoint name (required)")
	note := fs.String("note", "", "Optional note")
	
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *game == "" || *name == "" {
		return fmt.Errorf("both --game and --name are required")
	}

	fmt.Printf("Creating checkpoint...\n")
	
	checkpoint, err := c.service.CreateCheckpoint(*game, *name, *note)
	if err != nil {
		return fmt.Errorf("failed to create checkpoint: %w", err)
	}

	fmt.Printf("✓ Checkpoint created successfully\n")
	fmt.Printf("  ID:      %s\n", checkpoint.ID)
	fmt.Printf("  Name:    %s\n", checkpoint.Name)
	if checkpoint.Note != "" {
		fmt.Printf("  Note:    %s\n", checkpoint.Note)
	}
	fmt.Printf("  Created: %s\n", checkpoint.CreatedAt.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("  Hash:    %s\n", checkpoint.Hash[:16]+"...")
	
	return nil
}

// listCheckpoints handles the list command
func (c *CLI) listCheckpoints(args []string) error {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	game := fs.String("game", "", "Game ID or name (required)")
	
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *game == "" {
		return fmt.Errorf("--game is required")
	}

	checkpoints, err := c.service.ListCheckpoints(*game)
	if err != nil {
		return fmt.Errorf("failed to list checkpoints: %w", err)
	}

	if len(checkpoints) == 0 {
		fmt.Printf("No checkpoints found for game: %s\n", *game)
		return nil
	}

	fmt.Printf("Checkpoints for %s (%d):\n\n", *game, len(checkpoints))
	
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tCREATED\tNOTE")
	fmt.Fprintln(w, "──\t────\t───────\t────")
	
	for _, cp := range checkpoints {
		// Format created time
		created := cp.CreatedAt.Local().Format("2006-01-02 15:04")
		
		// Shorten ID for display
		shortID := cp.ID
		if len(shortID) > 8 {
			shortID = shortID[:8]
		}
		
		// Truncate note if too long
		note := cp.Note
		if len(note) > 40 {
			note = note[:37] + "..."
		}
		
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", shortID, cp.Name, created, note)
	}
	
	w.Flush()
	return nil
}

// restoreCheckpoint handles the restore command
func (c *CLI) restoreCheckpoint(args []string) error {
	fs := flag.NewFlagSet("restore", flag.ExitOnError)
	checkpoint := fs.String("checkpoint", "", "Checkpoint ID (required)")
	
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *checkpoint == "" {
		return fmt.Errorf("--checkpoint is required")
	}

	// Get checkpoint details
	cp, err := c.service.GetCheckpoint(*checkpoint)
	if err != nil {
		return fmt.Errorf("failed to get checkpoint: %w", err)
	}

	fmt.Printf("Restoring checkpoint...\n")
	fmt.Printf("  Name:    %s\n", cp.Name)
	fmt.Printf("  Created: %s\n", cp.CreatedAt.Format("2006-01-02 15:04:05 MST"))
	fmt.Printf("\n")

	if err := c.service.RestoreCheckpoint(*checkpoint); err != nil {
		return fmt.Errorf("failed to restore checkpoint: %w", err)
	}

	fmt.Printf("✓ Checkpoint restored successfully\n")
	
	return nil
}

// deleteCheckpoint handles the delete command
func (c *CLI) deleteCheckpoint(args []string) error {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	checkpoint := fs.String("checkpoint", "", "Checkpoint ID (required)")
	
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *checkpoint == "" {
		return fmt.Errorf("--checkpoint is required")
	}

	if err := c.service.DeleteCheckpoint(*checkpoint); err != nil {
		return fmt.Errorf("failed to delete checkpoint: %w", err)
	}

	fmt.Printf("✓ Checkpoint deleted successfully\n")
	
	return nil
}

// printUsage prints usage information
func (c *CLI) printUsage() {
	fmt.Printf(`GameKeep v%s - Game Save Manager (CLI)

USAGE:
    gamekeep <command> [options]

COMMANDS:
    add-game      Register a new game
    list-games    List all registered games
    checkpoint    Create a checkpoint for a game
    list          List checkpoints for a game
    restore       Restore a checkpoint
    delete        Delete a checkpoint
    version       Show version information
    help          Show this help message

EXAMPLES:
    # Register a game
    gamekeep add-game --name "The Witcher 3" --path "C:/Users/You/Documents/The Witcher 3"

    # Create a checkpoint
    gamekeep checkpoint --game witcher3 --name "Before Boss Fight" --note "Level 25, fire build"

    # List checkpoints
    gamekeep list --game witcher3

    # Restore a checkpoint
    gamekeep restore --checkpoint abc12345

    # Delete a checkpoint
    gamekeep delete --checkpoint abc12345

    # List all games
    gamekeep list-games

NOTE: For GUI interface, run 'gamekeep-gui' instead.

For more information, visit: https://github.com/adrielfilipedesign/gamekeep
`, version)
}

func main() {
	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to get home directory: %v\n", err)
		os.Exit(1)
	}

	// Setup directories
	basePath := filepath.Join(homeDir, ".gamekeep")
	configDir := filepath.Join(basePath, "config")
	vaultDir := filepath.Join(basePath, "vault")

	// Initialize storage
	store, err := storage.NewJSONStore(configDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to initialize storage: %v\n", err)
		os.Exit(1)
	}

	// Initialize vault manager
	vaultMgr, err := vault.NewManager(vaultDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to initialize vault: %v\n", err)
		os.Exit(1)
	}

	// Initialize service
	service := core.NewService(store, vaultMgr)

	// Initialize CLI
	cli := NewCLI(service)

	// Run
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
