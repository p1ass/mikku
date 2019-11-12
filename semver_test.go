package mikku

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_bumpVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tag     string
		typ     bumpType
		want    string
		wantErr bool
	}{
		{
			name:    "major bump",
			tag:     "v1.2.3",
			typ:     major,
			want:    "v2.0.0",
			wantErr: false,
		},
		{
			name:    "minor bump",
			tag:     "v1.2.3",
			typ:     minor,
			want:    "v1.3.0",
			wantErr: false,
		},
		{
			name:    "patch bump",
			tag:     "v1.2.3",
			typ:     patch,
			want:    "v1.2.4",
			wantErr: false,
		},
		{
			name:    "not semver tag",
			tag:     "v1.2.a",
			typ:     patch,
			want:    "",
			wantErr: true,
		},
		{
			name:    "no prefix tag",
			tag:     "1.2.3",
			typ:     patch,
			want:    "v1.2.4",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bumpVersion(tt.tag, tt.typ)
			if (err != nil) != tt.wantErr {
				t.Errorf("bumpVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("bumpVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createSemanticVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		versions []int
		want     string
	}{
		{
			name:     "v1.0.0",
			versions: []int{1, 0, 0},
			want:     "v1.0.0",
		},
		{
			name:     "v0.2.0",
			versions: []int{0, 2, 0},
			want:     "v0.2.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createSemanticVersion(tt.versions); got != tt.want {
				t.Errorf("createSemanticVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_strsToInts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		strs    []string
		want    []int
		wantErr bool
	}{
		{
			name:    "1.0.0",
			strs:    []string{"1", "0", "0"},
			want:    []int{1, 0, 0},
			wantErr: false,
		},
		{
			name:    "invalid string",
			strs:    []string{"1", "0", "invalid"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := strsToInts(tt.strs)
			if (err != nil) != tt.wantErr {
				t.Errorf("strsToInts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want) {
				t.Errorf("strsToInts() = %v, want %v", got, tt.want)
			}
		})
	}
}
