package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Window sizes
var (
	MainWindowSize      = fyne.NewSize(1000, 700)
	DialogSize          = fyne.NewSize(500, 400)
	SmallDialogSize     = fyne.NewSize(400, 200)
	CheckpointDialogSize = fyne.NewSize(600, 450)
)

// Custom theme colors
var (
	PrimaryColor   = color.NRGBA{R: 66, G: 135, B: 245, A: 255}  // Blue
	SuccessColor   = color.NRGBA{R: 52, G: 199, B: 89, A: 255}   // Green
	WarningColor   = color.NRGBA{R: 255, G: 159, B: 10, A: 255}  // Orange
	DangerColor    = color.NRGBA{R: 255, G: 69, B: 58, A: 255}   // Red
	BackgroundColor = color.NRGBA{R: 28, G: 28, B: 30, A: 255}   // Dark gray
)

// GameKeepTheme is a custom dark theme
type GameKeepTheme struct{}

func (t GameKeepTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return PrimaryColor
	case theme.ColorNameBackground:
		return BackgroundColor
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (t GameKeepTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t GameKeepTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t GameKeepTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

// Icons and emojis for UI
const (
	IconGame       = "🎮"
	IconCheckpoint = "💾"
	IconRestore    = "⏮️"
	IconDelete     = "🗑️"
	IconAdd        = "➕"
	IconFolder     = "📁"
	IconSuccess    = "✅"
	IconWarning    = "⚠️"
	IconInfo       = "ℹ️"
)
