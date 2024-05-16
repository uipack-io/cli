package base

import "github.com/charmbracelet/lipgloss"

// Style definitions.
var (
	// Title.

	TitleStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Border(lipgloss.RoundedBorder())

	// List.

	CheckMark = lipgloss.NewStyle().SetString("✓").
			Foreground(lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}).
			PaddingRight(1).
			PaddingLeft(1).
			String()

	StartedMark = lipgloss.NewStyle().SetString("»").
			Foreground(lipgloss.AdaptiveColor{Light: "#FFC700", Dark: "#FFC700"}).
			PaddingTop(1).
			PaddingRight(1).
			PaddingLeft(1).
			String()

	SubinfoMark = lipgloss.NewStyle().SetString("■").
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			PaddingRight(1).
			PaddingLeft(4).
			String()

	ListStarted = func(s string) string {
		return StartedMark + lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(s)
	}

	ListDone = func(s string) string {
		return CheckMark + lipgloss.NewStyle().
			Render(s)
	}

	ListSubinfo = func(s string) string {
		return SubinfoMark + lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
			Render(s)
	}
)
