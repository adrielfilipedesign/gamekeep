package main

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/adrielfilipedesign/gamekeep/internal/core"
	"github.com/adrielfilipedesign/gamekeep/internal/storage"
	"github.com/adrielfilipedesign/gamekeep/internal/vault"
	"github.com/adrielfilipedesign/gamekeep/ui"
)

const (
	appID      = "com.gamekeep.app"
	appVersion = "1.0.0"
)

func main() {
	// Initialize Fyne app
	a := app.NewWithID(appID)
	
	// Setup directories
	homeDir, err := os.UserHomeDir()
	if err != nil {
		showErrorDialog(a, "Startup Error", fmt.Sprintf("Failed to get home directory: %v", err))
		return
	}

	basePath := filepath.Join(homeDir, ".gamekeep")
	configDir := filepath.Join(basePath, "config")
	vaultDir := filepath.Join(basePath, "vault")

	// Initialize storage
	store, err := storage.NewJSONStore(configDir)
	if err != nil {
		showErrorDialog(a, "Startup Error", fmt.Sprintf("Failed to initialize storage: %v", err))
		return
	}

	// Initialize vault manager
	vaultMgr, err := vault.NewManager(vaultDir)
	if err != nil {
		showErrorDialog(a, "Startup Error", fmt.Sprintf("Failed to initialize vault: %v", err))
		return
	}

	// Initialize service
	service := core.NewService(store, vaultMgr)

	// Create main window
	mainWindow := a.NewWindow("GameKeep - Save Manager")
	mainWindow.Resize(ui.MainWindowSize)
	mainWindow.SetMaster()

	// Create main UI
	mainUI := ui.NewMainUI(mainWindow, service)
	mainWindow.SetContent(mainUI.Build())

	// Show and run
	mainWindow.ShowAndRun()
}

func showErrorDialog(a fyne.App, title, message string) {
	w := a.NewWindow(title)
	w.SetContent(container.NewVBox(
		widget.NewLabel(message),
		widget.NewButton("Exit", func() {
			os.Exit(1)
		}),
	))
	w.Resize(fyne.NewSize(400, 150))
	w.ShowAndRun()
}
