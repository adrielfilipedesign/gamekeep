package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/adrielfilipedesign/gamekeep/internal/core"
	"github.com/adrielfilipedesign/gamekeep/internal/models"
)

// MainUI is the main application UI controller
type MainUI struct {
	window          fyne.Window
	service         *core.Service
	gamesView       *GamesView
	checkpointsView *CheckpointsView
	currentGame     *models.Game
}

// NewMainUI creates a new main UI controller
func NewMainUI(window fyne.Window, service *core.Service) *MainUI {
	ui := &MainUI{
		window:  window,
		service: service,
	}

	ui.gamesView = NewGamesView(ui)
	ui.checkpointsView = NewCheckpointsView(ui)

	return ui
}

// Build constructs the main UI layout
func (m *MainUI) Build() fyne.CanvasObject {
	// Header
	header := m.buildHeader()

	// Main content area - split view
	leftPanel := m.buildLeftPanel()
	rightPanel := m.buildRightPanel()

	split := container.NewHSplit(leftPanel, rightPanel)
	split.Offset = 0.35 // 35% for games list, 65% for checkpoints

	// Main container
	content := container.NewBorder(
		header, // top
		nil,    // bottom
		nil,    // left
		nil,    // right
		split,  // center
	)

	return content
}

// buildHeader creates the top header bar
func (m *MainUI) buildHeader() fyne.CanvasObject {
	title := widget.NewLabelWithStyle(
		"🎮 GameKeep",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true},
	)

	subtitle := widget.NewLabel("Save State Manager")
	subtitle.TextStyle.Italic = true

	titleBox := container.NewVBox(title, subtitle)

	// Action buttons
	refreshBtn := widget.NewButton("🔄 Refresh", func() {
		m.RefreshAll()
	})

	aboutBtn := widget.NewButton("ℹ️ About", func() {
		m.showAboutDialog()
	})

	buttons := container.NewHBox(
		layout.NewSpacer(),
		refreshBtn,
		aboutBtn,
	)

	header := container.NewBorder(
		nil,
		nil,
		titleBox,
		buttons,
		nil,
	)

	return container.NewPadded(header)
}

// buildLeftPanel creates the games list panel
func (m *MainUI) buildLeftPanel() fyne.CanvasObject {
	return container.NewBorder(
		widget.NewLabelWithStyle(IconGame+" Games", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		m.gamesView.Build(),
	)
}

// buildRightPanel creates the checkpoints panel
func (m *MainUI) buildRightPanel() fyne.CanvasObject {
	title := widget.NewLabelWithStyle(IconCheckpoint+" Checkpoints", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	return container.NewBorder(
		title,
		nil,
		nil,
		nil,
		m.checkpointsView.Build(),
	)
}

// SelectGame handles game selection
func (m *MainUI) SelectGame(game *models.Game) {
	m.currentGame = game
	m.checkpointsView.LoadCheckpoints(game)
}

// RefreshAll refreshes all views
func (m *MainUI) RefreshAll() {
	m.gamesView.Refresh()
	if m.currentGame != nil {
		m.checkpointsView.LoadCheckpoints(m.currentGame)
	}
}

// GetService returns the core service
func (m *MainUI) GetService() *core.Service {
	return m.service
}

// GetWindow returns the main window
func (m *MainUI) GetWindow() fyne.Window {
	return m.window
}

// showAboutDialog shows the about dialog
func (m *MainUI) showAboutDialog() {
	content := container.NewVBox(
		widget.NewLabelWithStyle("GameKeep v1.0.0", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
		widget.NewLabel("Game Save State Manager"),
		widget.NewLabel(""),
		widget.NewLabel("Features:"),
		widget.NewLabel("• Create and manage save checkpoints"),
		widget.NewLabel("• Restore previous save states"),
		widget.NewLabel("• SHA256 integrity verification"),
		widget.NewLabel("• Compressed storage"),
		widget.NewLabel(""),
		widget.NewLabel("Built with Go + Fyne"),
		widget.NewLabel(""),
		widget.NewButton("Close", func() {}),
	)

	dialog := NewCustomDialog("About GameKeep", content, m.window)
	dialog.Show()
}
