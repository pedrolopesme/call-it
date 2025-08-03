package call

import (
	"testing"
)

func TestParseCurlCommand(t *testing.T) {
	tests := []struct {
		name        string
		curlCommand string
		wantURL     string
		wantMethod  string
		wantError   bool
	}{
		{
			name:        "Simple GET request",
			curlCommand: "curl https://httpbin.org/get",
			wantURL:     "https://httpbin.org/get",
			wantMethod:  "GET",
			wantError:   false,
		},
		{
			name:        "POST request with JSON",
			curlCommand: `curl -X POST https://httpbin.org/post -H "Content-Type: application/json" -d '{"key": "value"}'`,
			wantURL:     "https://httpbin.org/post",
			wantMethod:  "POST",
			wantError:   false,
		},
		{
			name:        "GET with headers",
			curlCommand: `curl -H "Authorization: Bearer token123" https://httpbin.org/get`,
			wantURL:     "https://httpbin.org/get",
			wantMethod:  "GET",
			wantError:   false,
		},
		{
			name:        "Invalid curl command",
			curlCommand: "not-curl https://example.com",
			wantError:   true,
		},
		{
			name:        "Empty command",
			curlCommand: "",
			wantError:   true,
		},
		{
			name: "Multiline curl command",
			curlCommand: `curl 'https://httpbin.org/post' \
  -H 'Content-Type: application/json' \
  -X POST \
  -d '{"test": "data"}'`,
			wantURL:    "https://httpbin.org/post",
			wantMethod: "POST",
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ParseCurlCommand(tt.curlCommand)
			
			if tt.wantError {
				if err == nil {
					t.Errorf("ParseCurlCommand() expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("ParseCurlCommand() error = %v, want nil", err)
				return
			}
			
			if config.URL != tt.wantURL {
				t.Errorf("ParseCurlCommand() URL = %v, want %v", config.URL, tt.wantURL)
			}
			
			if config.Method != tt.wantMethod {
				t.Errorf("ParseCurlCommand() Method = %v, want %v", config.Method, tt.wantMethod)
			}
		})
	}
}

func TestValidateCurlCommand(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		wantError   bool
	}{
		{
			name:      "Valid curl command",
			command:   "curl https://httpbin.org/get",
			wantError: false,
		},
		{
			name:      "Valid curl with options",
			command:   `curl -X POST -H "Content-Type: application/json" https://httpbin.org/post`,
			wantError: false,
		},
		{
			name:      "Invalid - doesn't start with curl",
			command:   "wget https://example.com",
			wantError: true,
		},
		{
			name:      "Empty command",
			command:   "",
			wantError: true,
		},
		{
			name:      "Whitespace only",
			command:   "   ",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCurlCommand(tt.command)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateCurlCommand() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}