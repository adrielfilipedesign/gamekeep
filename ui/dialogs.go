package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// ShowError displays an error dialog
func ShowError(window fyne.Window, message string, err error) {
	fullMessage := message
	if err != nil {
		fullMessage = fmt.Sprintf("%s\n\nError: %v", message, err)
	}

	dialog.ShowError(fmt.Errorf("%s", fullMessage), window)
}

// ShowSuccess displays a success notification
func ShowSuccess(window fyne.Window, message string) {
	content := container.NewVBox(
		widget.NewLabel(IconSuccess+" "+message),
	)

	d := dialog.NewCustom("Success", "OK", content, window)
	d.Resize(SmallDialogSize)
	d.Show()
}

// ShowInfo displays an info notification
func ShowInfo(window fyne.Window, message string) {
	content := container.NewVBox(
		widget.NewLabel(IconInfo+" "+message),
	)

	d := dialog.NewCustom("Information", "OK", content, window)
	d.Resize(SmallDialogSize)
	d.Show()
}

// ShowWarning displays a warning notification
func ShowWarning(window fyne.Window, message string) {
	content := container.NewVBox(
		widget.NewLabel(IconWarning+" "+message),
	)

	d := dialog.NewCustom("Warning", "OK", content, window)
	d.Resize(SmallDialogSize)
	d.Show()
}

// CustomDialog is a reusable custom dialog
type CustomDialog struct {
	dialog  *dialog.CustomDialog
	content fyne.CanvasObject
}

// NewCustomDialog creates a new custom dialog
func NewCustomDialog(title string, content fyne.CanvasObject, window fyne.Window) *CustomDialog {
	d := dialog.NewCustom(title, "Close", content, window)

	return &CustomDialog{
		dialog:  d,
		content: content,
	}
}

// Show displays the dialog
func (d *CustomDialog) Show() {
	d.dialog.Show()
}

// Hide hides the dialog
func (d *CustomDialog) Hide() {
	d.dialog.Hide()
}

// Resize resizes the dialog
func (d *CustomDialog) Resize(size fyne.Size) {
	d.dialog.Resize(size)
}
