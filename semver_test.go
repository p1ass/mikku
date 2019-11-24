package mikku

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_validSemver(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ver  string
		want bool
	}{
		{
			name: "v1.2.3",
			ver:  "v1.2.3",
			want: true,
		},
		{
			name: "v0.0.0",
			ver:  "v0.0.0",
			want: true,
		},
		{
			name: "v.2.3",
			ver:  "v.3.3",
			want: false,
		},
		{
			name: "v1..3",
			ver:  "v1..3",
			want: false,
		},
		{
			name: "v1.2.",
			ver:  "v1.2.",
			want: false,
		},
		{
			name: "1.2.3",
			ver:  "1.2.3",
			want: false,
		},
		{
			name: "v1.2.a",
			ver:  "1.2.a",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validSemver(tt.ver); got != tt.want {
				t.Errorf("validSemver() = %v, want %v", got, tt.want)
			}
		})
	}
}

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

func Test_determineNewTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		currentTag string
		typORVer   string
		want       string
		wantErr    bool
	}{
		{
			name:       "patch and valid tag",
			currentTag: "v1.0.0",
			typORVer:   "patch",
			want:       "v1.0.1",
			wantErr:    false,
		},
		{
			name:       "patch and invalid tag",
			currentTag: "1.0.0",
			typORVer:   "patch",
			want:       "",
			wantErr:    true,
		},
		{
			name:       "specify version and valid tag",
			currentTag: "",
			typORVer:   "v1.0.0",
			want:       "v1.0.0",
			wantErr:    false,
		},
		{
			name:       "specify version and invalid tag",
			currentTag: "",
			typORVer:   "1.0.0",
			want:       "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := determineNewTag(tt.currentTag, tt.typORVer)
			if (err != nil) != tt.wantErr {
				t.Errorf("determineNewTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("determineNewTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
