package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/adrielfilipedesign/gamekeep/internal/models"
)

// CheckpointsView handles the checkpoints display
type CheckpointsView struct {
	mainUI      *MainUI
	currentGame *models.Game
	list        *widget.List
	checkpoints []models.Checkpoint
	container   *fyne.Container
	emptyLabel  *widget.Label
}

// NewCheckpointsView creates a new checkpoints view
func NewCheckpointsView(mainUI *MainUI) *CheckpointsView {
	return &CheckpointsView{
		mainUI:      mainUI,
		checkpoints: []models.Checkpoint{},
	}
}

// Build creates the checkpoints view UI
func (v *CheckpointsView) Build() fyne.CanvasObject {
	// Empty state label
	v.emptyLabel = widget.NewLabel("← Select a game to view checkpoints")
	v.emptyLabel.Alignment = fyne.TextAlignCenter

	// Create list
	v.list = widget.NewList(
		func() int {
			return len(v.checkpoints)
		},
		func() fyne.CanvasObject {
			return v.createCheckpointCard()
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(v.checkpoints) {
				return
			}
			v.updateCheckpointCard(obj, v.checkpoints[id])
		},
	)

	// Buttons
	createBtn := widget.NewButton(IconAdd+" Create Checkpoint", func() {
		if v.currentGame == nil {
			ShowInfo(v.mainUI.GetWindow(), "Please select a game first")
			return
		}
		v.showCreateCheckpointDialog()
	})

	createBtn.Importance = widget.HighImportance

	buttons := container.NewHBox(createBtn)

	// Container with conditional content
	v.container = container.NewBorder(
		nil,
		buttons,
		nil,
		nil,
		v.emptyLabel,
	)

	return v.container
}

// createCheckpointCard creates a card template for a checkpoint
func (v *CheckpointsView) createCheckpointCard() fyne.CanvasObject {
	nameLabel := widget.NewLabelWithStyle("", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	dateLabel := widget.NewLabel("")
	noteLabel := widget.NewLabel("")
	noteLabel.Wrapping = fyne.TextWrapWord

	info := container.NewVBox(
		nameLabel,
		dateLabel,
		noteLabel,
	)

	restoreBtn := widget.NewButton(IconRestore+" Restore", func() {})
	restoreBtn.Importance = widget.HighImportance

	deleteBtn := widget.NewButton(IconDelete, func() {})
	deleteBtn.Importance = widget.DangerImportance

	actions := container.NewHBox(
		restoreBtn,
		deleteBtn,
	)

	card := container.NewBorder(
		nil,
		nil,
		nil,
		actions,
		info,
	)

	return container.NewPadded(card)
}

// updateCheckpointCard updates a card with checkpoint data
func (v *CheckpointsView) updateCheckpointCard(obj fyne.CanvasObject, cp models.Checkpoint) {
	card := obj.(*fyne.Container).Objects[0].(*fyne.Container)
	info := card.Objects[0].(*fyne.Container)
	actions := card.Objects[1].(*fyne.Container)

	// Update info
	nameLabel := info.Objects[0].(*widget.Label)
	dateLabel := info.Objects[1].(*widget.Label)
	noteLabel := info.Objects[2].(*widget.Label)

	nameLabel.SetText(IconCheckpoint + " " + cp.Name)
	dateLabel.SetText(fmt.Sprintf("Created: %s", cp.CreatedAt.Local().Format("2006-01-02 15:04")))

	if cp.Note != "" {
		noteLabel.SetText(fmt.Sprintf("Note: %s", cp.Note))
		noteLabel.Show()
	} else {
		noteLabel.Hide()
	}

	// Update buttons
	restoreBtn := actions.Objects[0].(*widget.Button)
	deleteBtn := actions.Objects[1].(*widget.Button)

	restoreBtn.OnTapped = func() {
		v.confirmRestore(cp)
	}

	deleteBtn.OnTapped = func() {
		v.confirmDelete(cp)
	}
}

// LoadCheckpoints loads checkpoints for a game
func (v *CheckpointsView) LoadCheckpoints(game *models.Game) {
	v.currentGame = game

	checkpoints, err := v.mainUI.GetService().ListCheckpoints(game.ID)
	if err != nil {
		ShowError(v.mainUI.GetWindow(), "Failed to load checkpoints", err)
		return
	}

	v.checkpoints = checkpoints

	// Update UI
	if len(v.checkpoints) > 0 {
		v.container.Objects[0] = v.list
	} else {
		emptyMsg := widget.NewLabel(fmt.Sprintf("No checkpoints yet for %s\nClick 'Create Checkpoint' to create one", game.Name))
		emptyMsg.Alignment = fyne.TextAlignCenter
		v.container.Objects[0] = emptyMsg
	}

	v.container.Refresh()
	v.list.Refresh()
}

// showCreateCheckpointDialog shows dialog to create checkpoint
func (v *CheckpointsView) showCreateCheckpointDialog() {
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("e.g., Before Boss Fight")

	noteEntry := widget.NewMultiLineEntry()
	noteEntry.SetPlaceHolder("Optional notes...")
	noteEntry.SetMinRowsVisible(3)

	form := container.NewVBox(
		widget.NewLabel("Checkpoint Name:"),
		nameEntry,
		widget.NewLabel(""),
		widget.NewLabel("Notes (optional):"),
		noteEntry,
	)

	d := dialog.NewCustomConfirm(
		"Create Checkpoint",
		"Create",
		"Cancel",
		form,
		func(confirmed bool) {
			if !confirmed {
				return
			}

			name := nameEntry.Text
			if name == "" {
				ShowInfo(v.mainUI.GetWindow(), "Please enter a checkpoint name")
				return
			}

			// Show progress
			progress := dialog.NewProgressInfinite(
				"Creating Checkpoint",
				fmt.Sprintf("Backing up saves for %s...", v.currentGame.Name),
				v.mainUI.GetWindow(),
			)
			progress.Show()

			// Create checkpoint
			go func() {
				cp, err := v.mainUI.GetService().CreateCheckpoint(
					v.currentGame.ID,
					name,
					noteEntry.Text,
				)

				// Close progress
				progress.Hide()

				if err != nil {
					ShowError(v.mainUI.GetWindow(), "Failed to create checkpoint", err)
					return
				}

				ShowSuccess(v.mainUI.GetWindow(), fmt.Sprintf("Checkpoint '%s' created successfully!", cp.Name))
				v.LoadCheckpoints(v.currentGame)
			}()
		},
		v.mainUI.GetWindow(),
	)

	d.Resize(CheckpointDialogSize)
	d.Show()
}

// confirmRestore shows confirmation dialog for restore
func (v *CheckpointsView) confirmRestore(cp models.Checkpoint) {
	message := fmt.Sprintf(
		"Restore checkpoint '%s'?\n\nCreated: %s\n\n⚠️  WARNING: This will replace your current save files!",
		cp.Name,
		cp.CreatedAt.Local().Format("2006-01-02 15:04:05"),
	)

	d := dialog.NewConfirm(
		"Confirm Restore",
		message,
		func(confirmed bool) {
			if !confirmed {
				return
			}

			// Show progress
			progress := dialog.NewProgressInfinite(
				"Restoring Checkpoint",
				"Extracting save files...",
				v.mainUI.GetWindow(),
			)
			progress.Show()

			// Restore
			go func() {
				err := v.mainUI.GetService().RestoreCheckpoint(cp.ID)
				progress.Hide()

				if err != nil {
					ShowError(v.mainUI.GetWindow(), "Failed to restore checkpoint", err)
					return
				}

				ShowSuccess(v.mainUI.GetWindow(), fmt.Sprintf("Checkpoint '%s' restored successfully!", cp.Name))
			}()
		},
		v.mainUI.GetWindow(),
	)

	d.Show()
}

// confirmDelete shows confirmation dialog for delete
func (v *CheckpointsView) confirmDelete(cp models.Checkpoint) {
	message := fmt.Sprintf(
		"Delete checkpoint '%s'?\n\nCreated: %s\n\nThis action cannot be undone.",
		cp.Name,
		cp.CreatedAt.Local().Format("2006-01-02 15:04:05"),
	)

	d := dialog.NewConfirm(
		"Confirm Delete",
		message,
		func(confirmed bool) {
			if !confirmed {
				return
			}

			err := v.mainUI.GetService().DeleteCheckpoint(cp.ID)
			if err != nil {
				ShowError(v.mainUI.GetWindow(), "Failed to delete checkpoint", err)
				return
			}

			ShowSuccess(v.mainUI.GetWindow(), "Checkpoint deleted")
			v.LoadCheckpoints(v.currentGame)
		},
		v.mainUI.GetWindow(),
	)

	d.Show()
}
