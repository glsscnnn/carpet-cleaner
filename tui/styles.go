package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Tab bar
	ActiveTab = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 2).
			MarginRight(1)

	InactiveTab = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Background(lipgloss.Color("#333333")).
			Padding(0, 2).
			MarginRight(1)

	// Status bar
	StatusBar = lipgloss.NewStyle().
			Background(lipgloss.Color("#333333")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1)

	// Warning
	WarningBanner = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFCC00")).
			Background(lipgloss.Color("#442200")).
			Padding(0, 1).
			Bold(true)

	// Result markers
	SuccessMarker = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	ErrorMarker   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).Bold(true)

	// Modal
	ModalBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2)

	ModalOverlay = lipgloss.NewStyle().
			Background(lipgloss.Color("#000000"))

	// Checkbox
	CheckboxChecked   = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Bold(true)
	CheckboxUnchecked = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))

	// Visibility badge
	PrivateBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFCC00")).
			Bold(true)
	PublicBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00CC00"))

	// Summary
	SummaryLine = lipgloss.NewStyle().Bold(true)
)