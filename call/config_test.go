package call

import (
	"os"
	"reflect"
	"testing"
)

func Test_getConfig(t *testing.T) {
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
			twd:     "./_example",
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
