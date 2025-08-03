package call

import (
	"fmt"
	"strings"

	"github.com/474420502/gcurl"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ParseCurlCommand parses a curl command string and returns a Config struct
// that can be used with the existing call-it infrastructure
func ParseCurlCommand(curlCommand string) (*Config, error) {
	// Clean up the curl command - remove leading/trailing whitespace
	originalLength := len(curlCommand)
	curlCommand = strings.TrimSpace(curlCommand)
	
	// Handle multiline curl commands by processing line continuation characters
	curlCommand = preprocessMultilineCurl(curlCommand)
	
	// Keep track of processing for potential debugging
	_ = originalLength // Avoid unused variable warning
	
	// Parse the curl command using gcurl with fallback approaches
	parsedCurl, err := parseWithFallbacks(curlCommand)
	if err != nil {
		return nil, fmt.Errorf("failed to parse curl command: %w", err)
	}

	// Extract basic information
	config := &Config{
		Name:   "Parsed from cURL",
		URL:    parsedCurl.ParsedURL.String(),
		Method: strings.ToUpper(parsedCurl.Method),
	}

	// Extract headers
	if len(parsedCurl.Header) > 0 {
		config.Header = make(map[string][]string)
		for key, values := range parsedCurl.Header {
			config.Header[key] = values
		}
	}

	// Extract body content
	if parsedCurl.Body != nil && parsedCurl.Body.Content != nil {
		// Convert the body content to string
		if bodyStr, ok := parsedCurl.Body.Content.(string); ok {
			config.Body = bodyStr
		} else if bodyBytes, ok := parsedCurl.Body.Content.([]byte); ok {
			config.Body = string(bodyBytes)
		} else {
			// For other types, use the String() method
			config.Body = parsedCurl.Body.String()
		}
	}

	// Set default method if not specified
	if config.Method == "" {
		config.Method = "GET"
	}

	return config, nil
}

// ValidateCurlCommand checks if a string looks like a valid curl command
func ValidateCurlCommand(command string) error {
	command = strings.TrimSpace(command)
	
	// Basic validation - should start with "curl"
	if !strings.HasPrefix(strings.ToLower(command), "curl") {
		return fmt.Errorf("command must start with 'curl'")
	}
	
	// Handle multiline curl commands by processing line continuation characters
	command = preprocessMultilineCurl(command)
	
	// Try to parse it to see if it's valid
	_, err := gcurl.Parse(command)
	if err != nil {
		// Provide more specific error context
		if strings.Contains(err.Error(), "multiple URLs") {
			return fmt.Errorf("parsing error - this might be due to special characters in headers or URL. Try simplifying the command or using double quotes instead of single quotes: %w", err)
		}
		return fmt.Errorf("invalid curl command: %w", err)
	}
	
	return nil
}

// preprocessMultilineCurl handles multiline curl commands by properly processing
// line continuation characters and normalizing whitespace
func preprocessMultilineCurl(command string) string {
	// First, handle the case where textinput converted newlines to spaces
	// Look for patterns like "\ -H" and replace with " -H"
	command = handleTextInputProcessing(command)
	
	// Handle different types of line breaks and continuation characters
	lines := strings.Split(command, "\n")
	var processedLines []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Skip empty lines
		if line == "" {
			continue
		}
		
		// Handle line continuation character at the end
		if strings.HasSuffix(line, " \\") {
			// Remove the backslash and trailing space
			line = strings.TrimSuffix(line, " \\")
			line = strings.TrimSpace(line)
		} else if strings.HasSuffix(line, "\\") {
			// Remove the backslash
			line = strings.TrimSuffix(line, "\\")
			line = strings.TrimSpace(line)
		}
		
		// Add the line to processed lines
		if line != "" {
			processedLines = append(processedLines, line)
		}
	}
	
	// Join all lines with a single space
	result := strings.Join(processedLines, " ")
	
	// Clean up any multiple spaces
	for strings.Contains(result, "  ") {
		result = strings.ReplaceAll(result, "  ", " ")
	}
	
	return result
}

// handleTextInputProcessing fixes issues caused by textinput component
// converting newlines to spaces in multiline curl commands
func handleTextInputProcessing(command string) string {
	// Pattern: backslash followed by space and then a curl option
	// Replace "\ -H" with " -H", "\ -X" with " -X", etc.
	curlOptions := []string{"-H", "-X", "-d", "--data", "--header", "-A", "--user-agent", "-u", "--user", "-b", "--cookie", "-e", "--referer"}
	
	result := command
	for _, option := range curlOptions {
		pattern := "\\ " + option
		replacement := " " + option
		result = strings.ReplaceAll(result, pattern, replacement)
	}
	
	// Also handle the general case of "\ " at the end of arguments
	result = strings.ReplaceAll(result, "\\ ", " ")
	
	return result
}

// parseWithFallbacks tries multiple parsing approaches to handle edge cases
func parseWithFallbacks(command string) (*gcurl.CURL, error) {
	var lastErr error
	
	// First attempt: Direct parsing
	result, err := gcurl.Parse(command)
	if err == nil {
		return result, nil
	}
	lastErr = err
	
	// Second attempt: Convert single quotes to double quotes in headers
	if strings.Contains(err.Error(), "multiple URLs") || strings.Contains(err.Error(), "misplaced") {
		// Replace single quotes with double quotes only in header values
		modified := convertQuotesInHeaders(command)
		result, err = gcurl.Parse(modified)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}
	
	// Third attempt: Simplify complex headers
	simplified := simplifyComplexHeaders(command)
	result, err = gcurl.Parse(simplified)
	if err == nil {
		return result, nil
	}
	
	// Return the original error if all attempts failed
	return nil, lastErr
}

// convertQuotesInHeaders converts single quotes to double quotes in header values
func convertQuotesInHeaders(command string) string {
	// This is a simple approach - replace single quotes in -H arguments
	result := command
	
	// Find all -H 'header: value' patterns and convert to -H "header: value"
	parts := strings.Split(result, "-H '")
	if len(parts) > 1 {
		for i := 1; i < len(parts); i++ {
			// Find the closing single quote
			endQuote := strings.Index(parts[i], "'")
			if endQuote != -1 {
				// Replace the opening and closing single quotes with double quotes
				header := parts[i][:endQuote]
				rest := parts[i][endQuote+1:]
				parts[i] = fmt.Sprintf(`"%s"%s`, header, rest)
			}
		}
		result = strings.Join(parts, `-H `)
	}
	
	return result
}

// simplifyComplexHeaders removes problematic headers that might cause parsing issues
func simplifyComplexHeaders(command string) string {
	// Remove headers that commonly cause issues
	problematicHeaders := []string{
		`-H 'sec-ch-ua: "Not)A;Brand";v="8", "Chromium";v="138", "Google Chrome";v="138"'`,
		`-H 'priority: u=0, i'`,
	}
	
	result := command
	for _, header := range problematicHeaders {
		result = strings.ReplaceAll(result, header, "")
	}
	
	// Clean up multiple spaces
	for strings.Contains(result, "  ") {
		result = strings.ReplaceAll(result, "  ", " ")
	}
	
	return strings.TrimSpace(result)
}