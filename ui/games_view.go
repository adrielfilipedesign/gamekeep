package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/gamekeep/gamekeep/internal/models"
)

// GamesView handles the games list display
type GamesView struct {
	mainUI    *MainUI
	list      *widget.List
	games     []models.Game
	container *fyne.Container
}

// NewGamesView creates a new games view
func NewGamesView(mainUI *MainUI) *GamesView {
	v := &GamesView{
		mainUI: mainUI,
		games:  []models.Game{},
	}
	v.loadGames()
	return v
}

// Build creates the games view UI
func (v *GamesView) Build() fyne.CanvasObject {
	// Create list
	v.list = widget.NewList(
		func() int {
			return len(v.games)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Game Name"),
			)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(v.games) {
				return
			}
			game := v.games[id]
			obj.(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf("%s %s", IconGame, game.Name))
		},
	)

	v.list.OnSelected = func(id widget.ListItemID) {
		if id < len(v.games) {
			v.mainUI.SelectGame(&v.games[id])
		}
	}

	// Add game button
	addBtn := widget.NewButton(IconAdd+" Add Game", func() {
		v.showAddGameDialog()
	})

	// Container
	v.container = container.NewBorder(
		nil,
		addBtn,
		nil,
		nil,
		v.list,
	)

	return v.container
}

// loadGames loads games from service
func (v *GamesView) loadGames() {
	games, err := v.mainUI.GetService().ListGames()
	if err != nil {
		ShowError(v.mainUI.GetWindow(), "Failed to load games", err)
		return
	}
	v.games = games
}

// Refresh reloads and updates the list
func (v *GamesView) Refresh() {
	v.loadGames()
	if v.list != nil {
		v.list.Refresh()
	}
}

// showAddGameDialog shows the dialog to add a new game
func (v *GamesView) showAddGameDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("Game name...")

	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("Save directory path...")

	browseBtn := widget.NewButton(IconFolder+" Browse", func() {
		dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
			if err != nil || dir == nil {
				return
			}
			pathEntry.SetText(dir.Path())
		}, v.mainUI.GetWindow())
	})

	form := container.NewVBox(
		widget.NewLabel("Game Name:"),
		nameEntry,
		widget.NewLabel(""),
		widget.NewLabel("Save Directory:"),
		container.NewBorder(nil, nil, nil, browseBtn, pathEntry),
	)

	// Create dialog
	d := dialog.NewCustomConfirm(
		"Add New Game",
		"Add",
		"Cancel",
		form,
		func(confirmed bool) {
			if !confirmed {
				return
			}

			name := nameEntry.Text
			path := pathEntry.Text

			if name == "" || path == "" {
				ShowInfo(v.mainUI.GetWindow(), "Please fill in all fields")
				return
			}

			// Add game
			game, err := v.mainUI.GetService().AddGame(name, path)
			if err != nil {
				ShowError(v.mainUI.GetWindow(), "Failed to add game", err)
				return
			}

			ShowSuccess(v.mainUI.GetWindow(), fmt.Sprintf("Game '%s' added successfully!", game.Name))
			v.Refresh()
		},
		v.mainUI.GetWindow(),
	)

	d.Resize(DialogSize)
	d.Show()
}
