package mikku

import (
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "release: no arguments",
			args:    []string{"", "release"},
			wantErr: false,
		},
		{
			name:    "release: only one argument",
			args:    []string{"", "release", "mikku"},
			wantErr: true,
		},
		{
			name:    "pr: no arguments",
			args:    []string{"", "pr"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Run(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
