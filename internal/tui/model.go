package tui

import (
	"fmt"
	"net/url"
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
	CurlInputView
	LoadingView
	ResultsView
)

// Model represents the TUI model
type Model struct {
	state        ViewState
	urlInput     textinput.Model
	attemptsInput textinput.Model
	concurrentInput textinput.Model
	headersInput textinput.Model
	bodyInput    textinput.Model
	curlInput    textinput.Model
	httpMethods  []string
	selectedMethod int
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
	animationFrame  int
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

	headersInput := textinput.New()
	headersInput.Placeholder = "Content-Type:application/json,Authorization:Bearer token"
	headersInput.CharLimit = 512
	headersInput.Width = 60

	bodyInput := textinput.New()
	bodyInput.Placeholder = `{"key": "value", "message": "Hello World"}`
	bodyInput.CharLimit = 1024
	bodyInput.Width = 60

	curlInput := textinput.New()
	curlInput.Placeholder = `curl -X POST https://httpbin.org/post -H "Content-Type: application/json" -d '{"key": "value"}'`
	curlInput.CharLimit = 8192  // Increased limit for complex curl commands
	curlInput.Width = 100

	// Initialize spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(primaryColor)

	return Model{
		state:           InputView,
		urlInput:        urlInput,
		attemptsInput:   attemptsInput,
		concurrentInput: concurrentInput,
		headersInput:    headersInput,
		bodyInput:       bodyInput,
		curlInput:       curlInput,
		httpMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "TRACE", "CONNECT", "PATCH"},
		selectedMethod:  0, // Default to GET
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
		case CurlInputView:
			return m.updateCurlInputView(msg)
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
		m.currentProgress = m.totalProgress  // Set to 100%
		m.statusMessage = "Calls completed!"
		return m, nil

	case callErrorMsg:
		m.state = InputView
		m.error = msg.error
		m.urlInput.Focus()
		return m, textinput.Blink

	case progressTickMsg:
		// Only update progress if we're in loading state
		if m.state == LoadingView {
			// Increment animation frame for smooth animations
			m.animationFrame++
			
			// Don't let progress exceed total-1 to reserve 100% for actual completion
			if m.currentProgress < m.totalProgress-1 {
				m.currentProgress++
				m.statusMessage = fmt.Sprintf("Completed %d/%d requests", m.currentProgress, m.totalProgress)
			} else if m.currentProgress == m.totalProgress-1 {
				// At 99%, show "Finishing up..." and keep waiting for actual completion
				m.statusMessage = "Finishing up..."
			}
			// Always continue progress updates until actual completion
			return m, m.simulateProgress()
		}
		return m, nil
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
	case "ctrl+p":
		// Switch to curl paste mode
		m.state = CurlInputView
		m.error = ""
		m.curlInput.Focus()
		return m, textinput.Blink
	case "tab", "shift+tab", "up", "down":
		// Switch between inputs, skipping body input if method doesn't support it
		if msg.String() == "tab" || msg.String() == "down" {
			m.activeInput++
			// Skip body input (index 5) if method doesn't support body
			if m.activeInput == 5 && !m.methodSupportsBody() {
				m.activeInput = 0 // Wrap around to beginning
			} else {
				m.activeInput = m.activeInput % 6
			}
		} else {
			m.activeInput--
			if m.activeInput < 0 {
				// Go to body if supported, otherwise method selector
				if m.methodSupportsBody() {
					m.activeInput = 5
				} else {
					m.activeInput = 4
				}
			}
			// Skip body input (index 5) if method doesn't support body when going backwards
			if m.activeInput == 5 && !m.methodSupportsBody() {
				m.activeInput = 4
			}
		}
		
		// Update focus
		m.urlInput.Blur()
		m.attemptsInput.Blur()
		m.concurrentInput.Blur()
		m.headersInput.Blur()
		m.bodyInput.Blur()
		
		switch m.activeInput {
		case 0:
			m.urlInput.Focus()
		case 1:
			m.attemptsInput.Focus()
		case 2:
			m.concurrentInput.Focus()
		case 3:
			m.headersInput.Focus()
		case 4:
			// HTTP method selector - no focus needed
		case 5:
			m.bodyInput.Focus()
		}
		
		return m, textinput.Blink
	case "left", "right":
		// Handle HTTP method selection when active
		if m.activeInput == 4 {
			if msg.String() == "right" {
				m.selectedMethod = (m.selectedMethod + 1) % len(m.httpMethods)
			} else {
				m.selectedMethod = (m.selectedMethod - 1 + len(m.httpMethods)) % len(m.httpMethods)
			}
		}
	case "enter":
		// Start the call
		return m.startCall()
	default:
		// Pass all other key messages to the active input (not method selector)
		switch m.activeInput {
		case 0:
			m.urlInput, cmd = m.urlInput.Update(msg)
		case 1:
			m.attemptsInput, cmd = m.attemptsInput.Update(msg)
		case 2:
			m.concurrentInput, cmd = m.concurrentInput.Update(msg)
		case 3:
			m.headersInput, cmd = m.headersInput.Update(msg)
		case 4:
			// Method selector doesn't need text input
		case 5:
			m.bodyInput, cmd = m.bodyInput.Update(msg)
		}
		cmds = append(cmds, cmd)
	}
	
	return m, tea.Batch(cmds...)
}

// updateCurlInputView handles curl input view updates
func (m Model) updateCurlInputView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		// Go back to main input view
		m.state = InputView
		m.error = ""
		m.urlInput.Focus()
		return m, textinput.Blink
	case "enter":
		// Parse curl command and populate fields
		return m.parseCurlAndContinue()
	default:
		// Pass key to curl input
		m.curlInput, cmd = m.curlInput.Update(msg)
	}
	
	return m, cmd
}

// parseCurlAndContinue parses the curl command and populates the form fields
func (m Model) parseCurlAndContinue() (Model, tea.Cmd) {
	curlCommand := strings.TrimSpace(m.curlInput.Value())
	if curlCommand == "" {
		m.error = "Please enter a curl command"
		return m, nil
	}

	// Process the curl command

	// Validate and parse curl command
	if err := call.ValidateCurlCommand(curlCommand); err != nil {
		m.error = fmt.Sprintf("Validation failed: %v", err)
		return m, nil
	}

	config, err := call.ParseCurlCommand(curlCommand)
	if err != nil {
		m.error = fmt.Sprintf("Parsing failed: %v", err)
		return m, nil
	}

	// Populate form fields from parsed config
	m.urlInput.SetValue(config.URL)
	
	// Set HTTP method
	for i, method := range m.httpMethods {
		if method == config.Method {
			m.selectedMethod = i
			break
		}
	}

	// Set headers
	if config.Header != nil {
		var headerPairs []string
		for key, values := range config.Header {
			for _, value := range values {
				headerPairs = append(headerPairs, fmt.Sprintf("%s:%s", key, value))
			}
		}
		m.headersInput.SetValue(strings.Join(headerPairs, ","))
	}

	// Set body
	if config.Body != "" {
		m.bodyInput.SetValue(config.Body)
	}

	// Switch back to input view
	m.state = InputView
	m.error = ""
	m.urlInput.Focus()
	
	return m, textinput.Blink
}

// startCall validates inputs and starts the HTTP calls
func (m Model) startCall() (Model, tea.Cmd) {
	// Clear previous error
	m.error = ""
	
	// Validate URL
	urlString := strings.TrimSpace(m.urlInput.Value())
	if urlString == "" {
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
	
	// Parse headers
	headers := make(map[string][]string)
	headersStr := strings.TrimSpace(m.headersInput.Value())
	if headersStr != "" {
		headerPairs := strings.Split(headersStr, ",")
		for _, pair := range headerPairs {
			parts := strings.SplitN(strings.TrimSpace(pair), ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				if key != "" && value != "" {
					headers[key] = []string{value}
				}
			}
		}
	}
	
	// Parse body (only for methods that support it)
	var body string
	if m.methodSupportsBody() {
		body = strings.TrimSpace(m.bodyInput.Value())
	}
	
	// Build call configuration with HTTP method, headers, and body
	config := call.Config{
		Name:               "TUI Request",
		Method:             m.httpMethods[m.selectedMethod],
		URL:                urlString,
		Attempts:           attempts,
		ConcurrentAttempts: concurrent,
		Header:             headers,
		Body:               body,
	}
	
	// Validate the config
	if err := config.CheckDefaults(); err != nil {
		m.error = err.Error()
		return m, nil
	}
	
	// Create ConcurrentCall with config
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		m.error = err.Error()
		return m, nil
	}
	
	callConfig := call.ConcurrentCall{
		URL:                parsedURL,
		Attempts:           attempts,
		ConcurrentAttempts: concurrent,
	}
	
	// Set the config using the public setter method
	callConfig.SetConfig(config)
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
	return tea.Batch(
		m.runCalls(),
		m.simulateProgress(),
	)
}

// runCalls executes the HTTP calls and returns the result
func (m Model) runCalls() tea.Cmd {
	return func() tea.Msg {
		results := m.callConfig.MakeIt()
		return callCompleteMsg{results: &results}
	}
}

// simulateProgress sends periodic progress updates
func (m Model) simulateProgress() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(time.Time) tea.Msg {
		return progressTickMsg{}
	})
}

// progressTickMsg triggers a progress update
type progressTickMsg struct{}


// View implements tea.Model
func (m Model) View() string {
	switch m.state {
	case InputView:
		return m.renderInputView()
	case CurlInputView:
		return m.renderCurlInputView()
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
	
	// Headers Input
	headersLabel := "HTTP Headers (key:value,key:value):"
	if m.activeInput == 3 {
		headersLabel = "► " + headersLabel
		b.WriteString(focusedLabelStyle.Render(headersLabel))
	} else {
		b.WriteString(labelStyle.Render(headersLabel))
	}
	b.WriteString("\n")
	b.WriteString(m.headersInput.View())
	b.WriteString("\n\n")
	
	// HTTP Method Selector
	methodLabel := "HTTP Method:"
	if m.activeInput == 4 {
		methodLabel = "► " + methodLabel
		b.WriteString(focusedLabelStyle.Render(methodLabel))
	} else {
		b.WriteString(labelStyle.Render(methodLabel))
	}
	b.WriteString("\n")
	b.WriteString(m.renderMethodSelector())
	b.WriteString("\n\n")
	
	// Body Input (only show for methods that support body)
	if m.methodSupportsBody() {
		bodyLabel := "Request Body (JSON, XML, text):"
		if m.activeInput == 5 {
			bodyLabel = "► " + bodyLabel
			b.WriteString(focusedLabelStyle.Render(bodyLabel))
		} else {
			b.WriteString(labelStyle.Render(bodyLabel))
		}
		b.WriteString("\n")
		b.WriteString(m.bodyInput.View())
		b.WriteString("\n\n")
	}
	
	// Error message
	if m.error != "" {
		b.WriteString(StatusMessage(m.error, "error"))
		b.WriteString("\n\n")
	}
	
	// Instructions
	instructions := "Press Tab to navigate • Left/Right for method • Enter to start • Ctrl+P to paste curl • Ctrl+C to quit"
	b.WriteString(helpStyle.Render(instructions))
	
	return baseStyle.Render(b.String())
}

// renderCurlInputView renders the curl input screen
func (m Model) renderCurlInputView() string {
	var b strings.Builder
	
	// Header
	b.WriteString(RenderLogo())
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Paste cURL Command"))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(fmt.Sprintf("Version: %s", version.String())))
	b.WriteString("\n\n")
	
	// Form
	b.WriteString(titleStyle.Render("📋 Paste Your cURL Command"))
	b.WriteString("\n\n")
	
	// Curl Input
	b.WriteString(focusedLabelStyle.Render("► cURL Command:"))
	b.WriteString("\n")
	b.WriteString(m.curlInput.View())
	b.WriteString("\n\n")
	
	// Error message
	if m.error != "" {
		b.WriteString(StatusMessage(m.error, "error"))
		b.WriteString("\n\n")
	}
	
	// Instructions
	instructions := "Enter to parse and continue • Esc to go back • Ctrl+C to quit"
	b.WriteString(helpStyle.Render(instructions))
	
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
		b.WriteString(ProgressBarAnimated(m.currentProgress, m.totalProgress, m.animationFrame))
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
	
	// Execution time and method
	duration := m.endTime.Sub(m.startTime)
	method := m.httpMethods[m.selectedMethod]
	b.WriteString(StatusMessage(fmt.Sprintf("Completed %s request in %v", method, duration), "success"))
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

// methodSupportsBody returns true if the HTTP method supports request body
func (m Model) methodSupportsBody() bool {
	method := m.httpMethods[m.selectedMethod]
	switch method {
	case "POST", "PUT", "PATCH", "DELETE":
		return true
	default:
		return false
	}
}

// renderMethodSelector renders the HTTP method selector
func (m Model) renderMethodSelector() string {
	var methods []string
	
	for i, method := range m.httpMethods {
		if i == m.selectedMethod {
			// Highlight selected method
			if m.activeInput == 4 {
				// Active and selected
				methods = append(methods, activeButtonStyle.Render(method))
			} else {
				// Selected but not active
				methods = append(methods, buttonStyle.Render(method))
			}
		} else {
			// Not selected
			methodStyle := lipgloss.NewStyle().
				Foreground(mutedColor).
				Padding(0, 1).
				Margin(0, 1)
			methods = append(methods, methodStyle.Render(method))
		}
	}
	
	return lipgloss.JoinHorizontal(lipgloss.Left, methods...)
}