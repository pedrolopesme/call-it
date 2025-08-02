package tui

import (
	"fmt"
	"strings"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Purple-focused color scheme
	primaryColor   = lipgloss.Color("#8B5CF6")   // Bright purple
	secondaryColor = lipgloss.Color("#A855F7")   // Medium purple
	successColor   = lipgloss.Color("#10B981")   // Green (kept for status)
	errorColor     = lipgloss.Color("#EF4444")   // Red (kept for errors)
	warningColor   = lipgloss.Color("#F59E0B")   // Amber (kept for warnings)
	mutedColor     = lipgloss.Color("#9CA3AF")   // Cool gray
	bgColor        = lipgloss.Color("#1E1B3C")   // Dark purple background
	surfaceColor   = lipgloss.Color("#2D2A4A")   // Purple surface
	accentColor    = lipgloss.Color("#C084FC")   // Light purple accent
	highlightColor = lipgloss.Color("#DDD6FE")   // Very light purple

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// Title styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1).
			Margin(1, 0)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true).
			Margin(0, 0, 1, 0)

	// Input styles
	inputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(0, 1).
			Width(60)

	focusedInputStyle = inputStyle.
				BorderForeground(primaryColor).
				Foreground(highlightColor)

	// Label styles
	labelStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(false)

	focusedLabelStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	// Button styles
	buttonStyle = lipgloss.NewStyle().
			Foreground(highlightColor).
			Background(secondaryColor).
			Padding(0, 2).
			Margin(0, 1).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor)

	activeButtonStyle = buttonStyle.
				Background(primaryColor).
				BorderForeground(primaryColor)

	// Status styles
	successStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	// Progress styles
	progressBarStyle = lipgloss.NewStyle().
				Width(50).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(accentColor)

	// Table styles
	tableHeaderStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Bold(true).
				Align(lipgloss.Center).
				Padding(0, 1)

	tableCellStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Align(lipgloss.Left).
			Foreground(highlightColor)

	// Card styles
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(accentColor).
			Padding(1, 2).
			Margin(1, 0)

	// Help styles
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Margin(1, 0)

	// Logo style
	logoStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Align(lipgloss.Center).
			Margin(1, 0)
)

// RenderLogo renders the application logo with gradient effect
func RenderLogo() string {
	// Define gradient colors from light to dark purple
	gradientColors := []lipgloss.Color{
		lipgloss.Color("#DDD6FE"), // Very light purple
		lipgloss.Color("#C084FC"), // Light purple
		lipgloss.Color("#A855F7"), // Medium purple  
		lipgloss.Color("#8B5CF6"), // Bright purple
		lipgloss.Color("#7C3AED"), // Deep purple
		lipgloss.Color("#6D28D9"), // Darker purple
	}
	
	// Create gradient effect for ASCII art lines
	asciiLines := []string{
		" ██████╗  █████╗  ██╗      ██╗           ██╗ ████████╗",
		"██╔════╝ ██╔══██╗ ██║      ██║           ██║ ╚══██╔══╝",
		"██║      ███████║ ██║      ██║           ██║    ██║   ",
		"██║      ██╔══██║ ██║      ██║           ██║    ██║   ",
		"╚██████╗ ██║  ██║ ███████╗ ███████╗      ██║    ██║   ",
		" ╚═════╝ ╚═╝  ╚═╝ ╚══════╝ ╚══════╝      ╚═╝    ╚═╝   ",
	}
	
	var gradientAscii strings.Builder
	gradientAscii.WriteString("\n")
	
	for i, line := range asciiLines {
		colorIndex := i % len(gradientColors)
		lineStyle := lipgloss.NewStyle().
			Foreground(gradientColors[colorIndex]).
			Bold(true)
		gradientAscii.WriteString(lineStyle.Render(line))
		gradientAscii.WriteString("\n")
	}
	
	// Create gradient title text
	title := "✨ C A L L   I T ✨"
	var gradientTitle strings.Builder
	
	for i, char := range title {
		if char == ' ' {
			gradientTitle.WriteString(" ")
			continue
		}
		colorIndex := (i / 2) % len(gradientColors) // Slower color change
		charStyle := lipgloss.NewStyle().
			Foreground(gradientColors[colorIndex]).
			Bold(true)
		gradientTitle.WriteString(charStyle.Render(string(char)))
	}
	
	gradientAscii.WriteString("\n")
	gradientAscii.WriteString(lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(60).
		Render(gradientTitle.String()))
	gradientAscii.WriteString("\n")
	
	return logoStyle.Render(gradientAscii.String())
}

// StatusMessage renders a status message with appropriate styling
func StatusMessage(message string, msgType string) string {
	switch msgType {
	case "success":
		return successStyle.Render("✅ " + message)
	case "error":
		return errorStyle.Render("❌ " + message)
	case "warning":
		return warningStyle.Render("⚠️ " + message)
	default:
		return message
	}
}

// ProgressBar renders a progress bar
func ProgressBar(current, total int) string {
	if total == 0 {
		return ""
	}
	
	percentage := float64(current) / float64(total)
	filled := int(percentage * 40) // 40 chars width
	
	bar := ""
	for i := 0; i < 40; i++ {
		if i < filled {
			bar += lipgloss.NewStyle().Foreground(primaryColor).Render("█")
		} else {
			bar += lipgloss.NewStyle().Foreground(mutedColor).Render("░")
		}
	}
	
	return progressBarStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			bar,
			lipgloss.NewStyle().Margin(0, 1).Render(
				lipgloss.NewStyle().Foreground(secondaryColor).Bold(true).Render(
					lipgloss.PlaceHorizontal(10, lipgloss.Right, 
						fmt.Sprintf("%.1f%%", percentage*100),
					),
				),
			),
		),
	)
}

