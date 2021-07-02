package version_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thepwagner/action-update/version"
)

func TestSemverish(t *testing.T) {
	tests := map[string]string{
		"v1":         "v1",
		"1":          "v1",
		"1.0":        "v1.0",
		"1.0.0":      "v1.0.0",
		"1.0.0.0":    "v1.0.0-0",
		"1.0.0-beta": "v1.0.0-beta",
		"1.0.0.beta": "v1.0.0-beta",
	}
	for in, expected := range tests {
		t.Run(in, func(t *testing.T) {
			assert.Equal(t, expected, version.Semverish(in))
		})
	}
}

func TestSemverSort(t *testing.T) {
	tests := map[string]struct {
		in       []string
		expected []string
	}{
		"semver": {
			in:       []string{"v3", "v1", "v2"},
			expected: []string{"v3", "v2", "v1"},
		},
		"by specificity": {
			in:       []string{"v1.0", "v1", "v1.0.0"},
			expected: []string{"v1.0.0", "v1.0", "v1"},
		},
		"semver-ish": {
			in:       []string{"3.0.3", "1.0.1", "2.0.2"},
			expected: []string{"3.0.3", "2.0.2", "1.0.1"},
		},
	}

	for lbl, tc := range tests {
		t.Run(lbl, func(t *testing.T) {
			assert.Equal(t, tc.expected, version.SemverSort(tc.in))
		})
	}
}
