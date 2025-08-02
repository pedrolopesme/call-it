package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pedrolopesme/call-it/internal/call"
	"github.com/pedrolopesme/call-it/internal/version"
)

// ViewState represents the current view of the TUI
type ViewState int

const (
	InputView ViewState = iota
	LoadingView
	ResultsView
)

// Model represents the TUI model
type Model struct {
	state        ViewState
	urlInput     textinput.Model
	attemptsInput textinput.Model
	concurrentInput textinput.Model
	activeInput  int
	spinner      spinner.Model
	results      *call.Result
	callConfig   *call.ConcurrentCall
	error        string
	startTime    time.Time
	endTime      time.Time
	currentProgress int
	totalProgress   int
	statusMessage   string
	width          int
	height         int
}

// NewModel creates a new TUI model
func NewModel() Model {
	// Initialize text inputs
	urlInput := textinput.New()
	urlInput.Placeholder = "https://httpbin.org/status/200"
	urlInput.Focus()
	urlInput.CharLimit = 256
	urlInput.Width = 50

	attemptsInput := textinput.New()
	attemptsInput.Placeholder = "5"
	attemptsInput.CharLimit = 10
	attemptsInput.Width = 20

	concurrentInput := textinput.New()
	concurrentInput.Placeholder = "3"
	concurrentInput.CharLimit = 10
	concurrentInput.Width = 20

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)

	return Model{
		state:           InputView,
		urlInput:        urlInput,
		attemptsInput:   attemptsInput,
		concurrentInput: concurrentInput,
		activeInput:     0,
		spinner:         s,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case InputView:
			return m.updateInputView(msg)
		case LoadingView:
			if msg.String() == "ctrl+c" || msg.String() == "q" {
				return m, tea.Quit
			}
		case ResultsView:
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "r", "enter":
				// Reset to input view
				m.state = InputView
				m.error = ""
				m.results = nil
				m.urlInput.Focus()
				return m, textinput.Blink
			}
		}

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case callStartMsg:
		m.state = LoadingView
		m.startTime = time.Now()
		m.currentProgress = 0
		m.totalProgress = msg.total
		m.statusMessage = "Starting HTTP calls..."
		cmds = append(cmds, m.spinner.Tick)
		cmds = append(cmds, m.startCalls())
		return m, tea.Batch(cmds...)

	case callProgressMsg:
		m.currentProgress = msg.current
		m.statusMessage = fmt.Sprintf("Completed %d/%d requests", msg.current, m.totalProgress)
		return m, nil

	case callCompleteMsg:
		m.state = ResultsView
		m.endTime = time.Now()
		m.results = msg.results
		m.statusMessage = "Calls completed!"
		return m, nil

	case callErrorMsg:
		m.state = InputView
		m.error = msg.error
		m.urlInput.Focus()
		return m, textinput.Blink
	}


	return m, tea.Batch(cmds...)
}

// updateInputView handles input view updates
func (m Model) updateInputView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "tab", "shift+tab", "up", "down":
		// Switch between inputs
		if msg.String() == "tab" || msg.String() == "down" {
			m.activeInput = (m.activeInput + 1) % 3
		} else {
			m.activeInput = (m.activeInput - 1 + 3) % 3
		}
		
		// Update focus
		m.urlInput.Blur()
		m.attemptsInput.Blur()
		m.concurrentInput.Blur()
		
		switch m.activeInput {
		case 0:
			m.urlInput.Focus()
		case 1:
			m.attemptsInput.Focus()
		case 2:
			m.concurrentInput.Focus()
		}
		
		return m, textinput.Blink
	case "enter":
		// Start the call
		return m.startCall()
	default:
		// Pass all other key messages to the active input
		switch m.activeInput {
		case 0:
			m.urlInput, cmd = m.urlInput.Update(msg)
		case 1:
			m.attemptsInput, cmd = m.attemptsInput.Update(msg)
		case 2:
			m.concurrentInput, cmd = m.concurrentInput.Update(msg)
		}
		cmds = append(cmds, cmd)
	}
	
	return m, tea.Batch(cmds...)
}

// startCall validates inputs and starts the HTTP calls
func (m Model) startCall() (Model, tea.Cmd) {
	// Clear previous error
	m.error = ""
	
	// Validate URL
	url := strings.TrimSpace(m.urlInput.Value())
	if url == "" {
		m.error = "URL is required"
		return m, nil
	}
	
	// Parse attempts
	attemptsStr := strings.TrimSpace(m.attemptsInput.Value())
	if attemptsStr == "" {
		attemptsStr = "5"
	}
	attempts, err := strconv.Atoi(attemptsStr)
	if err != nil || attempts <= 0 {
		m.error = "Attempts must be a positive number"
		return m, nil
	}
	
	// Parse concurrent calls
	concurrentStr := strings.TrimSpace(m.concurrentInput.Value())
	if concurrentStr == "" {
		concurrentStr = "3"
	}
	concurrent, err := strconv.Atoi(concurrentStr)
	if err != nil || concurrent <= 0 {
		m.error = "Concurrent calls must be a positive number"
		return m, nil
	}
	
	// Build call configuration
	args := []string{url}
	callConfig, err := call.BuildCall(args, attempts, concurrent)
	if err != nil {
		m.error = err.Error()
		return m, nil
	}
	
	m.callConfig = &callConfig
	
	return m, func() tea.Msg {
		return callStartMsg{total: attempts}
	}
}

// Message types for async operations
type callStartMsg struct {
	total int
}

type callProgressMsg struct {
	current int
}

type callCompleteMsg struct {
	results *call.Result
}

type callErrorMsg struct {
	error string
}

// startCalls performs the actual HTTP calls
func (m Model) startCalls() tea.Cmd {
	return func() tea.Msg {
		// Simulate progress updates (in real implementation, you'd modify the call package)
		results := m.callConfig.MakeIt()
		return callCompleteMsg{results: &results}
	}
}

// View implements tea.Model
func (m Model) View() string {
	switch m.state {
	case InputView:
		return m.renderInputView()
	case LoadingView:
		return m.renderLoadingView()
	case ResultsView:
		return m.renderResultsView()
	default:
		return "Unknown view state"
	}
}

// renderInputView renders the input form
func (m Model) renderInputView() string {
	var b strings.Builder
	
	// Header
	b.WriteString(RenderLogo())
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("HTTP Load Testing Tool"))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(fmt.Sprintf("Version: %s", version.String())))
	b.WriteString("\n\n")
	
	// Form
	b.WriteString(titleStyle.Render("🚀 Configuration"))
	b.WriteString("\n\n")
	
	// URL Input
	urlLabel := "URL:"
	if m.activeInput == 0 {
		urlLabel = "► " + urlLabel
		b.WriteString(focusedLabelStyle.Render(urlLabel))
	} else {
		b.WriteString(labelStyle.Render(urlLabel))
	}
	b.WriteString("\n")
	b.WriteString(m.urlInput.View())
	b.WriteString("\n\n")
	
	// Attempts Input
	attemptsLabel := "Number of Attempts:"
	if m.activeInput == 1 {
		attemptsLabel = "► " + attemptsLabel
		b.WriteString(focusedLabelStyle.Render(attemptsLabel))
	} else {
		b.WriteString(labelStyle.Render(attemptsLabel))
	}
	b.WriteString("\n")
	b.WriteString(m.attemptsInput.View())
	b.WriteString("\n\n")
	
	// Concurrent Input
	concurrentLabel := "Concurrent Calls:"
	if m.activeInput == 2 {
		concurrentLabel = "► " + concurrentLabel
		b.WriteString(focusedLabelStyle.Render(concurrentLabel))
	} else {
		b.WriteString(labelStyle.Render(concurrentLabel))
	}
	b.WriteString("\n")
	b.WriteString(m.concurrentInput.View())
	b.WriteString("\n\n")
	
	// Error message
	if m.error != "" {
		b.WriteString(StatusMessage(m.error, "error"))
		b.WriteString("\n\n")
	}
	
	// Instructions
	b.WriteString(helpStyle.Render("Press Tab to navigate • Enter to start • Ctrl+C to quit"))
	
	return baseStyle.Render(b.String())
}

// renderLoadingView renders the loading screen
func (m Model) renderLoadingView() string {
	var b strings.Builder
	
	b.WriteString(titleStyle.Render("🔥 Running Load Test"))
	b.WriteString("\n\n")
	
	b.WriteString(fmt.Sprintf("%s %s", m.spinner.View(), m.statusMessage))
	b.WriteString("\n\n")
	
	if m.totalProgress > 0 {
		b.WriteString(ProgressBar(m.currentProgress, m.totalProgress))
		b.WriteString("\n\n")
	}
	
	b.WriteString(helpStyle.Render("Press Ctrl+C or q to quit"))
	
	return baseStyle.Render(b.String())
}

// renderResultsView renders the results
func (m Model) renderResultsView() string {
	var b strings.Builder
	
	b.WriteString(titleStyle.Render("📊 Results"))
	b.WriteString("\n\n")
	
	// Execution time
	duration := m.endTime.Sub(m.startTime)
	b.WriteString(StatusMessage(fmt.Sprintf("Completed in %v", duration), "success"))
	b.WriteString("\n\n")
	
	// Results table
	b.WriteString(cardStyle.Render(m.formatResults()))
	b.WriteString("\n\n")
	
	b.WriteString(helpStyle.Render("Press r or Enter to run again • Ctrl+C or q to quit"))
	
	return baseStyle.Render(b.String())
}

// formatResults formats the results into a readable table
func (m Model) formatResults() string {
	if m.results == nil {
		return "No results to display"
	}
	
	var b strings.Builder
	
	// Stats header
	b.WriteString(tableHeaderStyle.Render("Execution Stats"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 50))
	b.WriteString("\n")
	
	b.WriteString(fmt.Sprintf("Total Execution Time: %.2fs\n", m.results.GetTotalExecution()))
	b.WriteString(fmt.Sprintf("Average Execution Time: %.2fs\n", m.results.GetAvgExecution()))
	b.WriteString(fmt.Sprintf("Min Execution Time: %.2fs\n", m.results.GetMinExecution()))
	b.WriteString(fmt.Sprintf("Max Execution Time: %.2fs\n", m.results.GetMaxExecution()))
	b.WriteString("\n")
	
	// Status codes table
	b.WriteString(tableHeaderStyle.Render("Status Code"))
	b.WriteString("  ")
	b.WriteString(tableHeaderStyle.Render("Count"))
	b.WriteString("  ")
	b.WriteString(tableHeaderStyle.Render("Avg Time"))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 50))
	b.WriteString("\n")
	
	statusMap := m.results.GetStatus()
	if len(statusMap) == 0 {
		return "No status codes to display"
	}
	
	for status, benchmark := range statusMap {
		statusStr := fmt.Sprintf("%d", status)
		countStr := fmt.Sprintf("%d", benchmark.GetTotal())
		avgTimeStr := fmt.Sprintf("%.2fs", benchmark.GetExecution()/float64(benchmark.GetTotal()))
		
		// Color code based on HTTP status
		var statusStyle lipgloss.Style
		if status >= 200 && status < 300 {
			statusStyle = successStyle
		} else if status >= 400 {
			statusStyle = errorStyle
		} else {
			statusStyle = warningStyle
		}
		
		b.WriteString(statusStyle.Render(fmt.Sprintf("%-12s", statusStr)))
		b.WriteString(tableCellStyle.Render(fmt.Sprintf("%-8s", countStr)))
		b.WriteString(tableCellStyle.Render(avgTimeStr))
		b.WriteString("\n")
	}
	
	return b.String()
}