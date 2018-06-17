package call

import (
	"net/http"
	"os"
	"reflect"
	"testing"
)

func Test_config(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		wantC   []Config
		wantErr bool
		twd     string
	}{
		{
			name:    "should return err",
			wantErr: true,
		},
		{
			name:    "should read file",
			wantErr: false,
			twd:     "../_example",
			wantC: []Config{
				Config{
					Name:   "test request",
					Method: "GET",
					URL:    "http://www.globo.com",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.twd != "" {
				os.Chdir(tt.twd)
			}
			gotC, err := config()
			if (err != nil) != tt.wantErr {
				t.Errorf("config() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("config() = %v, want %v", gotC, tt.wantC)
			}
			os.Chdir(wd)
		})
	}
}

func TestConfig_checkDefaults(t *testing.T) {
	type fields struct {
		Name               string
		Method             string
		Attempts           int
		ConcurrentAttempts int
		URL                string
		Body               string
		Header             map[string][]string
		Host               string
		Form               string
		PostForm           map[string][]string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "config should have a name",
			fields:  fields{},
			wantErr: true,
		},
		{
			name:    "config method not allowed",
			fields:  fields{Name: "something", URL: "http://survivingmars.com", Method: "ASHE"},
			wantErr: true,
		},
		{
			name:    "invalid url",
			fields:  fields{Name: "something", URL: "survivingmars.com"},
			wantErr: true,
		},
		{
			name:    "config should pass",
			fields:  fields{Name: "something", URL: "http://survivingmars.com", Method: http.MethodGet},
			wantErr: false,
		},
		{
			name:    "empty url should not pass",
			fields:  fields{Name: "something", URL: "", Method: http.MethodGet},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Name:               tt.fields.Name,
				Method:             tt.fields.Method,
				Attempts:           tt.fields.Attempts,
				ConcurrentAttempts: tt.fields.ConcurrentAttempts,
				URL:                tt.fields.URL,
				Body:               tt.fields.Body,
				Header:             tt.fields.Header,
				Host:               tt.fields.Host,
				Form:               tt.fields.Form,
				PostForm:           tt.fields.PostForm,
			}
			if err := c.checkDefaults(); (err != nil) != tt.wantErr {
				t.Errorf("Config.checkDefaults() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
