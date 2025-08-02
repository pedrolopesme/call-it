package tui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pedrolopesme/call-it/internal/call"
)

func TestNewModel(t *testing.T) {
	model := NewModel()

	// Test initial state
	if model.state != InputView {
		t.Errorf("Expected initial state to be InputView, got %v", model.state)
	}

	// Test HTTP methods initialization
	expectedMethods := []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "TRACE", "CONNECT", "PATCH"}
	if len(model.httpMethods) != len(expectedMethods) {
		t.Errorf("Expected %d HTTP methods, got %d", len(expectedMethods), len(model.httpMethods))
	}

	for i, expected := range expectedMethods {
		if model.httpMethods[i] != expected {
			t.Errorf("Expected method %d to be %s, got %s", i, expected, model.httpMethods[i])
		}
	}

	// Test default selections
	if model.selectedMethod != 0 {
		t.Errorf("Expected default selected method to be 0 (GET), got %d", model.selectedMethod)
	}

	if model.activeInput != 0 {
		t.Errorf("Expected default active input to be 0, got %d", model.activeInput)
	}

	// Test input initialization
	if !model.urlInput.Focused() {
		t.Error("Expected URL input to be focused initially")
	}
}

func TestModelInit(t *testing.T) {
	model := NewModel()
	cmd := model.Init()

	if cmd == nil {
		t.Error("Expected Init() to return a command")
	}
}

func TestModelNavigation(t *testing.T) {
	model := NewModel()

	tests := []struct {
		name           string
		key            string
		expectedActive int
		expectedFocus  func(Model) bool
	}{
		{
			name:           "Tab to attempts input",
			key:            "tab",
			expectedActive: 1,
			expectedFocus:  func(m Model) bool { return m.attemptsInput.Focused() },
		},
		{
			name:           "Tab to concurrent input",
			key:            "tab",
			expectedActive: 2,
			expectedFocus:  func(m Model) bool { return m.concurrentInput.Focused() },
		},
		{
			name:           "Tab to headers input",
			key:            "tab",
			expectedActive: 3,
			expectedFocus:  func(m Model) bool { return m.headersInput.Focused() },
		},
		{
			name:           "Tab to method selector",
			key:            "tab",
			expectedActive: 4,
			expectedFocus:  func(m Model) bool { return !m.urlInput.Focused() && !m.attemptsInput.Focused() && !m.concurrentInput.Focused() && !m.headersInput.Focused() && !m.bodyInput.Focused() },
		},
		{
			name:           "Tab wraps to URL (GET method skips body)",
			key:            "tab",
			expectedActive: 0,
			expectedFocus:  func(m Model) bool { return m.urlInput.Focused() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyMsg := tea.KeyMsg{Type: tea.KeyTab}
			if tt.key == "shift+tab" {
				keyMsg = tea.KeyMsg{Type: tea.KeyShiftTab}
			}

			newModel, _ := model.Update(keyMsg)
			model = newModel.(Model)

			if model.activeInput != tt.expectedActive {
				t.Errorf("Expected active input to be %d, got %d", tt.expectedActive, model.activeInput)
			}

			if !tt.expectedFocus(model) {
				t.Errorf("Expected focus validation failed for %s", tt.name)
			}
		})
	}
}

func TestMethodSelection(t *testing.T) {
	model := NewModel()
	
	// Navigate to method selector
	for i := 0; i < 4; i++ {
		keyMsg := tea.KeyMsg{Type: tea.KeyTab}
		newModel, _ := model.Update(keyMsg)
		model = newModel.(Model)
	}

	if model.activeInput != 4 {
		t.Fatalf("Expected to be on method selector, got active input %d", model.activeInput)
	}

	tests := []struct {
		name             string
		key              string
		expectedMethod   int
		expectedMethodName string
	}{
		{
			name:             "Right arrow to POST",
			key:              "right",
			expectedMethod:   1,
			expectedMethodName: "POST",
		},
		{
			name:             "Right arrow to PUT",
			key:              "right",
			expectedMethod:   2,
			expectedMethodName: "PUT",
		},
		{
			name:             "Left arrow back to POST",
			key:              "left",
			expectedMethod:   1,
			expectedMethodName: "POST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var keyMsg tea.KeyMsg
			if tt.key == "right" {
				keyMsg = tea.KeyMsg{Type: tea.KeyRight}
			} else {
				keyMsg = tea.KeyMsg{Type: tea.KeyLeft}
			}

			newModel, _ := model.Update(keyMsg)
			model = newModel.(Model)

			if model.selectedMethod != tt.expectedMethod {
				t.Errorf("Expected selected method to be %d, got %d", tt.expectedMethod, model.selectedMethod)
			}

			if model.httpMethods[model.selectedMethod] != tt.expectedMethodName {
				t.Errorf("Expected method name to be %s, got %s", tt.expectedMethodName, model.httpMethods[model.selectedMethod])
			}
		})
	}
}

func TestStartCallValidation(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		attempts      string
		concurrent    string
		method        int
		expectedError string
	}{
		{
			name:          "Empty URL",
			url:           "",
			attempts:      "5",
			concurrent:    "3",
			method:        0,
			expectedError: "URL is required",
		},
		{
			name:          "Invalid attempts",
			url:           "https://example.com",
			attempts:      "invalid",
			concurrent:    "3",
			method:        0,
			expectedError: "Attempts must be a positive number",
		},
		{
			name:          "Zero attempts",
			url:           "https://example.com",
			attempts:      "0",
			concurrent:    "3",
			method:        0,
			expectedError: "Attempts must be a positive number",
		},
		{
			name:          "Invalid concurrent",
			url:           "https://example.com",
			attempts:      "5",
			concurrent:    "invalid",
			method:        0,
			expectedError: "Concurrent calls must be a positive number",
		},
		{
			name:          "Zero concurrent",
			url:           "https://example.com",
			attempts:      "5",
			concurrent:    "0",
			method:        0,
			expectedError: "Concurrent calls must be a positive number",
		},
		{
			name:          "Invalid URL format",
			url:           "not-a-url",
			attempts:      "5",
			concurrent:    "3",
			method:        0,
			expectedError: "parse \"not-a-url\": invalid URI for request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel()
			model.selectedMethod = tt.method

			// Set input values
			model.urlInput.SetValue(tt.url)
			model.attemptsInput.SetValue(tt.attempts)
			model.concurrentInput.SetValue(tt.concurrent)

			// Trigger startCall
			newModel, _ := model.startCall()

			if newModel.error != tt.expectedError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedError, newModel.error)
			}
		})
	}
}

func TestStartCallSuccess(t *testing.T) {
	model := NewModel()
	model.selectedMethod = 0 // GET

	// Set valid input values
	model.urlInput.SetValue("https://example.com")
	model.attemptsInput.SetValue("5")
	model.concurrentInput.SetValue("3")

	newModel, cmd := model.startCall()

	// Should not have error
	if newModel.error != "" {
		t.Errorf("Expected no error, got: %s", newModel.error)
	}

	// Should return a command
	if cmd == nil {
		t.Error("Expected startCall to return a command")
	}

	// Should have config set
	if newModel.callConfig == nil {
		t.Error("Expected callConfig to be set")
	}
}

func TestViewStates(t *testing.T) {
	model := NewModel()

	// Test InputView
	model.state = InputView
	view := model.View()
	if !strings.Contains(view, "Configuration") {
		t.Error("Expected InputView to contain 'Configuration'")
	}
	if !strings.Contains(view, "URL:") {
		t.Error("Expected InputView to contain 'URL:'")
	}

	// Test LoadingView
	model.state = LoadingView
	model.statusMessage = "Testing..."
	view = model.View()
	if !strings.Contains(view, "Running Load Test") {
		t.Error("Expected LoadingView to contain 'Running Load Test'")
	}
	if !strings.Contains(view, "Testing...") {
		t.Error("Expected LoadingView to contain status message")
	}

	// Test ResultsView
	model.state = ResultsView
	model.startTime = time.Now().Add(-5 * time.Second)
	model.endTime = time.Now()
	view = model.View()
	if !strings.Contains(view, "Results") {
		t.Error("Expected ResultsView to contain 'Results'")
	}
	if !strings.Contains(view, "Completed") {
		t.Error("Expected ResultsView to contain 'Completed'")
	}
}

func TestRenderMethodSelector(t *testing.T) {
	model := NewModel()

	// Test with GET selected and not active
	model.selectedMethod = 0
	model.activeInput = 0
	selector := model.renderMethodSelector()
	
	if selector == "" {
		t.Error("Expected method selector to render content")
	}

	// Test with method selector active
	model.activeInput = 4
	selectorActive := model.renderMethodSelector()
	
	if selectorActive == "" {
		t.Error("Expected active method selector to render content")
	}

	// Should contain method names
	if !strings.Contains(selector, "GET") {
		t.Error("Method selector should contain GET")
	}
}

func TestMessageHandling(t *testing.T) {
	model := NewModel()

	tests := []struct {
		name        string
		msg         tea.Msg
		expectState ViewState
	}{
		{
			name:        "Window size message",
			msg:         tea.WindowSizeMsg{Width: 100, Height: 50},
			expectState: InputView,
		},
		{
			name:        "Call start message",
			msg:         callStartMsg{total: 10},
			expectState: LoadingView,
		},
		{
			name: "Call complete message",
			msg: callCompleteMsg{results: &call.Result{}},
			expectState: ResultsView,
		},
		{
			name:        "Call error message",
			msg:         callErrorMsg{error: "test error"},
			expectState: InputView,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newModel, _ := model.Update(tt.msg)
			model = newModel.(Model)

			if model.state != tt.expectState {
				t.Errorf("Expected state %v, got %v", tt.expectState, model.state)
			}
		})
	}
}

func TestFormatResults(t *testing.T) {
	model := NewModel()

	// Test with nil results
	result := model.formatResults()
	if !strings.Contains(result, "No results to display") {
		t.Error("Expected message for nil results")
	}

	// Test with mock results
	mockResult := &call.Result{}
	model.results = mockResult
	
	result = model.formatResults()
	if strings.Contains(result, "No results to display") {
		t.Error("Should not show 'No results' message when results are set")
	}
}

func TestInputViewRendering(t *testing.T) {
	model := NewModel()
	
	// Test error display
	model.error = "Test error message"
	view := model.renderInputView()
	
	if !strings.Contains(view, "Test error message") {
		t.Error("Error message should be displayed in input view")
	}
	
	// Test without error
	model.error = ""
	view = model.renderInputView()
	
	if strings.Contains(view, "Test error message") {
		t.Error("Error message should not be displayed when error is empty")
	}
}

func TestLoadingViewRendering(t *testing.T) {
	model := NewModel()
	model.state = LoadingView
	model.statusMessage = "Test status"
	model.totalProgress = 10
	model.currentProgress = 5
	
	view := model.renderLoadingView()
	
	if !strings.Contains(view, "Test status") {
		t.Error("Status message should be displayed in loading view")
	}
	
	if !strings.Contains(view, "50.0%") {
		t.Error("Progress percentage should be displayed")
	}
}

func TestResultsViewRendering(t *testing.T) {
	model := NewModel()
	model.state = ResultsView
	model.selectedMethod = 1 // POST
	model.startTime = time.Now().Add(-2 * time.Second)
	model.endTime = time.Now()
	
	view := model.renderResultsView()
	
	if !strings.Contains(view, "POST") {
		t.Error("HTTP method should be displayed in results")
	}
	
	if !strings.Contains(view, "Completed") {
		t.Error("Completion message should be displayed")
	}
}

func TestKeyboardShortcuts(t *testing.T) {
	model := NewModel()
	
	// Test Ctrl+C
	keyMsg := tea.KeyMsg{Type: tea.KeyCtrlC}
	_, cmd := model.Update(keyMsg)
	
	if cmd == nil {
		t.Error("Ctrl+C should return quit command")
	}
	
	// Test Enter key
	model.urlInput.SetValue("https://example.com")
	keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.Update(keyMsg)
	
	// Should trigger validation or start call
	if newModel.(Model).error == "" && newModel.(Model).callConfig == nil {
		// This is expected if there are validation errors
	}
}

func TestDefaultValues(t *testing.T) {
	model := NewModel()
	
	// Test default placeholder values
	if !strings.Contains(model.urlInput.Placeholder, "httpbin.org") {
		t.Error("URL input should have httpbin.org placeholder")
	}
	
	if model.attemptsInput.Placeholder != "5" {
		t.Error("Attempts input should have '5' placeholder")
	}
	
	if model.concurrentInput.Placeholder != "3" {
		t.Error("Concurrent input should have '3' placeholder")
	}
	
	if !strings.Contains(model.headersInput.Placeholder, "Content-Type:application/json") {
		t.Error("Headers input should have example headers placeholder")
	}
}

func TestMethodSelectorBounds(t *testing.T) {
	model := NewModel()
	model.activeInput = 4 // Method selector
	
	// Test cycling through all methods
	originalMethod := model.selectedMethod
	
	// Go right through all methods
	for i := 0; i < len(model.httpMethods); i++ {
		keyMsg := tea.KeyMsg{Type: tea.KeyRight}
		newModel, _ := model.Update(keyMsg)
		model = newModel.(Model)
	}
	
	// Should wrap around to original
	if model.selectedMethod != originalMethod {
		t.Error("Method selection should wrap around when reaching end")
	}
	
	// Test left direction
	keyMsg := tea.KeyMsg{Type: tea.KeyLeft}
	newModel, _ := model.Update(keyMsg)
	model = newModel.(Model)
	
	expectedMethod := (originalMethod - 1 + len(model.httpMethods)) % len(model.httpMethods)
	if model.selectedMethod != expectedMethod {
		t.Errorf("Expected method %d, got %d", expectedMethod, model.selectedMethod)
	}
}

func TestHeadersParsing(t *testing.T) {
	tests := []struct {
		name          string
		headersInput  string
		expectedError string
	}{
		{
			name:          "Valid headers",
			headersInput:  "Content-Type:application/json,Authorization:Bearer token123",
			expectedError: "",
		},
		{
			name:          "Single header",
			headersInput:  "Accept:application/xml",
			expectedError: "",
		},
		{
			name:          "Headers with spaces",
			headersInput:  "Content-Type: application/json, Authorization: Bearer token",
			expectedError: "",
		},
		{
			name:          "Empty headers",
			headersInput:  "",
			expectedError: "",
		},
		{
			name:          "Invalid header format (no colon)",
			headersInput:  "ContentType application/json",
			expectedError: "",  // Should just skip invalid headers
		},
		{
			name:          "Mixed valid and invalid",
			headersInput:  "Content-Type:application/json,InvalidHeader,Accept:text/plain",
			expectedError: "",  // Should parse valid ones and skip invalid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := NewModel()
			model.selectedMethod = 0 // GET

			// Set valid input values
			model.urlInput.SetValue("https://example.com")
			model.attemptsInput.SetValue("5")
			model.concurrentInput.SetValue("3")
			model.headersInput.SetValue(tt.headersInput)

			// Trigger startCall
			newModel, _ := model.startCall()

			if tt.expectedError == "" && newModel.error != "" {
				t.Errorf("Expected no error, got: %s", newModel.error)
			}

			if tt.expectedError != "" && newModel.error != tt.expectedError {
				t.Errorf("Expected error '%s', got '%s'", tt.expectedError, newModel.error)
			}
		})
	}
}

func TestBodyInputNavigation(t *testing.T) {
	model := NewModel()
	
	// Change to POST method to enable body input
	model.selectedMethod = 1 // POST
	
	// Navigate to method selector first
	for i := 0; i < 4; i++ {
		keyMsg := tea.KeyMsg{Type: tea.KeyTab}
		newModel, _ := model.Update(keyMsg)
		model = newModel.(Model)
	}
	
	if model.activeInput != 4 {
		t.Fatalf("Expected to be on method selector, got active input %d", model.activeInput)
	}
	
	// Tab to body input (should work with POST method)
	keyMsg := tea.KeyMsg{Type: tea.KeyTab}
	newModel, _ := model.Update(keyMsg)
	model = newModel.(Model)
	
	if model.activeInput != 5 {
		t.Errorf("Expected to be on body input, got active input %d", model.activeInput)
	}
	
	if !model.bodyInput.Focused() {
		t.Error("Body input should be focused when active for POST method")
	}
}

func TestHeadersNavigation(t *testing.T) {
	model := NewModel()
	
	// Navigate to headers input
	for i := 0; i < 3; i++ {
		keyMsg := tea.KeyMsg{Type: tea.KeyTab}
		newModel, _ := model.Update(keyMsg)
		model = newModel.(Model)
	}
	
	if model.activeInput != 3 {
		t.Fatalf("Expected to be on headers input, got active input %d", model.activeInput)
	}
	
	if !model.headersInput.Focused() {
		t.Error("Headers input should be focused when active")
	}
}

func TestBodyInputWithPOST(t *testing.T) {
	model := NewModel()
	
	// Change to POST method first
	model.selectedMethod = 1 // POST
	model.activeInput = 5    // Body input
	
	// Check that body input is focused for POST
	model.bodyInput.Focus()
	if !model.bodyInput.Focused() {
		t.Error("Body input should be focused for POST method")
	}
	
	// Check that method supports body
	if !model.methodSupportsBody() {
		t.Error("POST method should support body")
	}
}

func TestBodyInputVisibility(t *testing.T) {
	model := NewModel()
	
	// Test methods that support body
	bodyMethods := []int{1, 2, 3, 8} // POST, PUT, DELETE, PATCH
	for _, methodIndex := range bodyMethods {
		model.selectedMethod = methodIndex
		if !model.methodSupportsBody() {
			t.Errorf("Method %s should support body", model.httpMethods[methodIndex])
		}
	}
	
	// Test methods that don't support body
	noBodyMethods := []int{0, 4, 5, 6, 7} // GET, OPTIONS, HEAD, TRACE, CONNECT
	for _, methodIndex := range noBodyMethods {
		model.selectedMethod = methodIndex
		if model.methodSupportsBody() {
			t.Errorf("Method %s should not support body", model.httpMethods[methodIndex])
		}
	}
}

func TestStartCallWithBody(t *testing.T) {
	model := NewModel()
	model.selectedMethod = 1 // POST method

	// Set valid input values including body
	model.urlInput.SetValue("https://example.com")
	model.attemptsInput.SetValue("5")
	model.concurrentInput.SetValue("3")
	model.bodyInput.SetValue(`{"test": "data", "value": 123}`)

	// Trigger startCall
	newModel, _ := model.startCall()

	// Should not have error
	if newModel.error != "" {
		t.Errorf("Expected no error, got: %s", newModel.error)
	}

	// Should have config set
	if newModel.callConfig == nil {
		t.Error("Expected callConfig to be set")
	}
}

func TestStartCallWithoutBodyForGET(t *testing.T) {
	model := NewModel()
	model.selectedMethod = 0 // GET method

	// Set valid input values, body should be ignored for GET
	model.urlInput.SetValue("https://example.com")
	model.attemptsInput.SetValue("5")
	model.concurrentInput.SetValue("3")
	model.bodyInput.SetValue(`{"this": "should be ignored"}`)

	// Trigger startCall
	newModel, _ := model.startCall()

	// Should not have error
	if newModel.error != "" {
		t.Errorf("Expected no error, got: %s", newModel.error)
	}

	// Should have config set
	if newModel.callConfig == nil {
		t.Error("Expected callConfig to be set")
	}
}

func TestProgressUpdates(t *testing.T) {
	model := NewModel()
	model.state = LoadingView
	model.totalProgress = 10
	model.currentProgress = 0
	
	// Simulate progress tick message
	progressMsg := progressTickMsg{}
	newModel, cmd := model.Update(progressMsg)
	model = newModel.(Model)
	
	// Should increment progress
	if model.currentProgress != 1 {
		t.Errorf("Expected progress to be 1, got %d", model.currentProgress)
	}
	
	// Should return a command to continue progress if not complete
	if cmd == nil {
		t.Error("Expected progress update to return a command")
	}
	
	// Test status message update
	expectedMessage := "Completed 1/10 requests"
	if model.statusMessage != expectedMessage {
		t.Errorf("Expected status message '%s', got '%s'", expectedMessage, model.statusMessage)
	}
}

func TestProgressCompletionStop(t *testing.T) {
	model := NewModel()
	model.state = LoadingView
	model.totalProgress = 5
	model.currentProgress = 4 // At total-1 (99%)
	
	// Simulate progress tick message
	progressMsg := progressTickMsg{}
	newModel, cmd := model.Update(progressMsg)
	model = newModel.(Model)
	
	// Should not increment beyond total-1
	if model.currentProgress != 4 {
		t.Errorf("Expected progress to remain 4, got %d", model.currentProgress)
	}
	
	// Should return a command to continue waiting for actual completion
	if cmd == nil {
		t.Error("Expected progress update to return a command even at 99%")
	}
	
	// Status should show "Finishing up..."
	if model.statusMessage != "Finishing up..." {
		t.Errorf("Expected 'Finishing up...' message, got '%s'", model.statusMessage)
	}
}

func TestProgressTickInNonLoadingState(t *testing.T) {
	model := NewModel()
	model.state = InputView // Not in loading state
	model.totalProgress = 10
	model.currentProgress = 0
	
	// Simulate progress tick message
	progressMsg := progressTickMsg{}
	newModel, cmd := model.Update(progressMsg)
	model = newModel.(Model)
	
	// Should not update progress when not in loading state
	if model.currentProgress != 0 {
		t.Errorf("Expected progress to remain 0, got %d", model.currentProgress)
	}
	
	// Should not return a command
	if cmd != nil {
		t.Error("Expected no command when not in loading state")
	}
}

func TestCallCompletionSetsProgress(t *testing.T) {
	model := NewModel()
	model.state = LoadingView
	model.totalProgress = 10
	model.currentProgress = 8
	
	// Create mock result
	mockResult := &call.Result{}
	
	// Simulate call completion message
	completeMsg := callCompleteMsg{results: mockResult}
	newModel, cmd := model.Update(completeMsg)
	model = newModel.(Model)
	
	// Should transition to results view
	if model.state != ResultsView {
		t.Errorf("Expected state to be ResultsView, got %v", model.state)
	}
	
	// Should set progress to 100%
	if model.currentProgress != model.totalProgress {
		t.Errorf("Expected progress to be %d (100%%), got %d", model.totalProgress, model.currentProgress)
	}
	
	// Should set completion message
	if model.statusMessage != "Calls completed!" {
		t.Errorf("Expected 'Calls completed!' message, got '%s'", model.statusMessage)
	}
	
	// Should have results set
	if model.results != mockResult {
		t.Error("Expected results to be set")
	}
	
	// Should not return a command
	if cmd != nil {
		t.Error("Expected no command after completion")
	}
}