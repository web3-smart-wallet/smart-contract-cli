package constant

import "github.com/charmbracelet/lipgloss"

var (
    // Error styling
    ErrorStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FF0000")).
        Bold(true).
        Padding(0, 1)

    // Success styling
    SuccessStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#00FF00")).
        Bold(true).
        Padding(0, 1)

    // Info styling
    InfoStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#00FFFF")).
        Bold(true).
        Padding(0, 1)
) 