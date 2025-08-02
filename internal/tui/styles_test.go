package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestColorDefinitions(t *testing.T) {
	tests := []struct {
		name  string
		color lipgloss.Color
	}{
		{"primaryColor", primaryColor},
		{"secondaryColor", secondaryColor},
		{"successColor", successColor},
		{"errorColor", errorColor},
		{"warningColor", warningColor},
		{"mutedColor", mutedColor},
		{"bgColor", bgColor},
		{"surfaceColor", surfaceColor},
		{"accentColor", accentColor},
		{"highlightColor", highlightColor},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.color) == "" {
				t.Errorf("Color %s should not be empty", tt.name)
			}
			
			// Check if it's a valid hex color (starts with #)
			colorStr := string(tt.color)
			if !strings.HasPrefix(colorStr, "#") {
				t.Errorf("Color %s should be a hex color, got: %s", tt.name, colorStr)
			}
			
			// Check if it's the right length (7 characters including #)
			if len(colorStr) != 7 {
				t.Errorf("Color %s should be 7 characters long, got: %d", tt.name, len(colorStr))
			}
		})
	}
}

func TestStyleDefinitions(t *testing.T) {
	tests := []struct {
		name  string
		style lipgloss.Style
	}{
		{"baseStyle", baseStyle},
		{"titleStyle", titleStyle},
		{"subtitleStyle", subtitleStyle},
		{"inputStyle", inputStyle},
		{"focusedInputStyle", focusedInputStyle},
		{"labelStyle", labelStyle},
		{"focusedLabelStyle", focusedLabelStyle},
		{"buttonStyle", buttonStyle},
		{"activeButtonStyle", activeButtonStyle},
		{"successStyle", successStyle},
		{"errorStyle", errorStyle},
		{"warningStyle", warningStyle},
		{"progressBarStyle", progressBarStyle},
		{"tableHeaderStyle", tableHeaderStyle},
		{"tableCellStyle", tableCellStyle},
		{"cardStyle", cardStyle},
		{"helpStyle", helpStyle},
		{"logoStyle", logoStyle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that styles can render content without panicking
			rendered := tt.style.Render("test")
			if rendered == "" {
				t.Errorf("Style %s should render content", tt.name)
			}
		})
	}
}

func TestRenderLogo(t *testing.T) {
	logo := RenderLogo()
	
	// Should not be empty
	if logo == "" {
		t.Error("RenderLogo should return non-empty string")
	}
	
	// Should contain ASCII art characters
	if !strings.Contains(logo, "██") {
		t.Error("Logo should contain ASCII art block characters")
	}
	
	// Should contain the sparkle emojis and text
	if !strings.Contains(logo, "✨") {
		t.Error("Logo should contain sparkle emojis")
	}
	
	// Should contain gradient text
	if !strings.Contains(logo, "C A L L   I T") {
		t.Error("Logo should contain spaced title text")
	}
}

func TestStatusMessage(t *testing.T) {
	tests := []struct {
		name        string
		message     string
		msgType     string
		expectEmoji string
	}{
		{
			name:        "Success message",
			message:     "Operation completed",
			msgType:     "success",
			expectEmoji: "✅",
		},
		{
			name:        "Error message",
			message:     "Something went wrong",
			msgType:     "error",
			expectEmoji: "❌",
		},
		{
			name:        "Warning message",
			message:     "Be careful",
			msgType:     "warning",
			expectEmoji: "⚠️",
		},
		{
			name:        "Default message",
			message:     "Just a message",
			msgType:     "info",
			expectEmoji: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StatusMessage(tt.message, tt.msgType)
			
			// Should contain the original message
			if !strings.Contains(result, tt.message) {
				t.Errorf("Result should contain original message: %s", tt.message)
			}
			
			// Should contain expected emoji (if any)
			if tt.expectEmoji != "" && !strings.Contains(result, tt.expectEmoji) {
				t.Errorf("Result should contain emoji: %s", tt.expectEmoji)
			}
			
			// Default case should return unchanged message
			if tt.msgType == "info" && result != tt.message {
				t.Errorf("Default case should return unchanged message, got: %s", result)
			}
		})
	}
}

func TestProgressBar(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		total    int
		expected string
	}{
		{
			name:     "Zero total",
			current:  5,
			total:    0,
			expected: "",
		},
		{
			name:     "Zero progress",
			current:  0,
			total:    10,
			expected: "0.0%",
		},
		{
			name:     "Half progress",
			current:  5,
			total:    10,
			expected: "50.0%",
		},
		{
			name:     "Full progress",
			current:  10,
			total:    10,
			expected: "100.0%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProgressBar(tt.current, tt.total)
			
			if tt.expected == "" && result != "" {
				t.Errorf("Expected empty result for zero total, got: %s", result)
			}
			
			if tt.expected != "" {
				if !strings.Contains(result, tt.expected) {
					t.Errorf("Expected result to contain %s, got: %s", tt.expected, result)
				}
				
				// Should contain progress bar characters (new colorful version uses different chars)
				if !strings.Contains(result, "█") && !strings.Contains(result, "▒") && !strings.Contains(result, "▓") {
					t.Error("Progress bar should contain progress characters")
				}
			}
		})
	}
}

func TestProgressBarEdgeCases(t *testing.T) {
	// Test negative values
	result := ProgressBar(-1, 10)
	if !strings.Contains(result, "0.0%") {
		t.Error("Negative current should be treated as 0%")
	}
	
	// Test current > total
	result = ProgressBar(15, 10)
	if !strings.Contains(result, "150.0%") {
		t.Error("Current > total should show >100%")
	}
}

func TestProgressBarAnimated(t *testing.T) {
	// Test animated version
	result := ProgressBarAnimated(5, 10, 0)
	if result == "" {
		t.Error("Animated progress bar should return content")
	}
	
	// Should contain percentage
	if !strings.Contains(result, "50.0%") {
		t.Error("Animated progress bar should contain percentage")
	}
	
	// Should contain progress characters
	if !strings.Contains(result, "█") && !strings.Contains(result, "▒") {
		t.Error("Animated progress bar should contain progress characters")
	}
	
	// Test animation frames produce different results
	result1 := ProgressBarAnimated(9, 10, 0)
	result2 := ProgressBarAnimated(9, 10, 1)
	if result1 == result2 {
		t.Error("Different animation frames should produce different results")
	}
	
	// Test completion shows checkmark
	resultComplete := ProgressBarAnimated(10, 10, 0)
	if !strings.Contains(resultComplete, "✅") {
		t.Error("Completed progress should show checkmark")
	}
}

func TestProgressBarColorGradient(t *testing.T) {
	// Test that progress bars at different percentages have different visual appearance
	// (We can't easily test colors, but we can test structure)
	
	tests := []struct {
		name     string
		current  int
		total    int
		expected string
	}{
		{"Low progress", 1, 10, "10.0%"},
		{"Medium progress", 5, 10, "50.0%"},
		{"High progress", 9, 10, "90.0%"},
		{"Complete", 10, 10, "100.0%"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProgressBarAnimated(tt.current, tt.total, 0)
			if !strings.Contains(result, tt.expected) {
				t.Errorf("Expected progress bar to contain %s, got: %s", tt.expected, result)
			}
		})
	}
}

func TestProgressBarLayoutStability(t *testing.T) {
	// Test that progress bar maintains consistent layout across animation frames
	current, total := 5, 10
	
	// Collect results from multiple animation frames
	var results []string
	for frame := 0; frame < 10; frame++ {
		result := ProgressBarAnimated(current, total, frame)
		results = append(results, result)
	}
	
	// All results should have exactly the same length (single-line layout)
	baseLength := len(results[0])
	for i, result := range results {
		if len(result) != baseLength {
			t.Errorf("Frame %d has different length: %d vs base %d", 
				i, len(result), baseLength)
		}
		
		// All should contain the same percentage
		if !strings.Contains(result, "50.0%") {
			t.Errorf("Frame %d missing expected percentage", i)
		}
		
		// Should be single line (no newlines)
		if strings.Contains(result, "\n") {
			t.Errorf("Frame %d contains newline characters (should be single line)", i)
		}
	}
}


func TestStyleConsistency(t *testing.T) {
	// Test that styles can render without errors
	regularInput := inputStyle.Render("test")
	focusedInput := focusedInputStyle.Render("test")
	
	if regularInput == "" || focusedInput == "" {
		t.Error("Input styles should render content")
	}
	
	regularLabel := labelStyle.Render("test")
	focusedLabel := focusedLabelStyle.Render("test")
	
	if regularLabel == "" || focusedLabel == "" {
		t.Error("Label styles should render content")
	}
	
	regularButton := buttonStyle.Render("test")
	activeButton := activeButtonStyle.Render("test")
	
	if regularButton == "" || activeButton == "" {
		t.Error("Button styles should render content")
	}
}

func TestPurpleColorScheme(t *testing.T) {
	// Test that primary colors are purple-ish
	purpleColors := []lipgloss.Color{primaryColor, secondaryColor, accentColor, highlightColor}
	
	for i, color := range purpleColors {
		colorStr := string(color)
		// Purple colors typically have high blue and red components
		// These specific colors should be in purple range
		expectedPurples := []string{"#8B5CF6", "#A855F7", "#C084FC", "#DDD6FE"}
		if colorStr != expectedPurples[i] {
			t.Errorf("Expected purple color %s, got %s", expectedPurples[i], colorStr)
		}
	}
}

func TestGradientColors(t *testing.T) {
	// Test that gradient colors in RenderLogo are properly ordered
	logo := RenderLogo()
	
	// Should be non-empty and contain colorized content
	if logo == "" {
		t.Error("Gradient logo should not be empty")
	}
	
	// The gradient should create a visually appealing effect
	// We can't easily test the actual color codes without complex parsing,
	// but we can ensure the basic structure is there
	if !strings.Contains(logo, "C A L L") || !strings.Contains(logo, "I T") {
		t.Error("Gradient logo should contain the spaced text content")
	}
}

func TestStyleRenderingIntegration(t *testing.T) {
	// Test that all styles can work together without conflicts
	content := titleStyle.Render("Title") + 
		subtitleStyle.Render("Subtitle") +
		inputStyle.Render("Input") +
		helpStyle.Render("Help")
	
	if content == "" {
		t.Error("Combined styles should render content")
	}
	
	// Test card with table content
	tableContent := tableHeaderStyle.Render("Header") + tableCellStyle.Render("Cell")
	cardContent := cardStyle.Render(tableContent)
	
	if cardContent == "" {
		t.Error("Nested styles should render content")
	}
}