package styles

import (
	"github.com/fortissolucoescontato-bit/kortex/internal/system/brand"
	"github.com/charmbracelet/lipgloss"
)

// Kortex Elite color palette (Carbon & Electric Blue).
var (
	ColorBase     = lipgloss.Color("#0f172a") // Deep Slate
	ColorSurface  = lipgloss.Color("#1e293b") // Slate
	ColorOverlay  = lipgloss.Color("#475569") // Muted Slate
	ColorText     = lipgloss.Color("#f8fafc") // Cloud White
	ColorSubtext  = lipgloss.Color("#94a3b8") // Light Slate
	ColorLavender = lipgloss.Color("#3b82f6") // Electric Blue
	ColorGreen    = lipgloss.Color("#10b981") // Emerald
	ColorPeach    = lipgloss.Color("#f59e0b") // Amber
	ColorRed      = lipgloss.Color("#ef4444") // Crimson
	ColorBlue     = lipgloss.Color("#2563eb") // Royal Blue
	ColorMauve    = lipgloss.Color("#6366f1") // Indigo
	ColorYellow   = lipgloss.Color("#fbbf24") // Gold
	ColorTeal     = lipgloss.Color("#14b8a6") // Cyan
)

// Cursor is the prefix used for the currently focused item.
const Cursor = "⚡ "

// Tagline returns the welcome screen tagline.
func Tagline() string {
	return brand.Tagline
}

// Pre-built reusable styles.
var (
	TitleStyle = lipgloss.NewStyle().
			Foreground(ColorLavender).
			Bold(true)

	HeadingStyle = lipgloss.NewStyle().
			Foreground(ColorMauve).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	SubtextStyle = lipgloss.NewStyle().
			Foreground(ColorSubtext)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(ColorLavender).
			Bold(true)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(ColorText)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorGreen)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorRed)

	WarningStyle = lipgloss.NewStyle().
			Foreground(ColorYellow)

	FrameStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(ColorLavender).
			Padding(1, 2)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(ColorOverlay).
			Padding(0, 1)

	ProgressFilled = lipgloss.NewStyle().
			Foreground(ColorGreen)

	ProgressEmpty = lipgloss.NewStyle().
			Foreground(ColorOverlay)

	PercentStyle = lipgloss.NewStyle().
			Foreground(ColorPeach).
			Bold(true)
)
