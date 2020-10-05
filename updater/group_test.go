package updater_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thepwagner/action-update/updater"
)

func TestGroup_InRange(t *testing.T) {
	cases := map[string]struct {
		included []string
		excluded []string
	}{
		"> v1": {
			included: []string{"v2", "v1.0.1"},
			excluded: []string{"v1", "v1.0.0", "v0.1"},
		},
		">= v1": {
			included: []string{"v2", "v1.0.1", "v1", "v1.0.0"},
			excluded: []string{"v0.1"},
		},
		"< v1": {
			included: []string{"v0.1"},
			excluded: []string{"v2", "v1.0.1", "v1", "v1.0.0"},
		},
		"<= v1": {
			included: []string{"v0.1", "v1", "v1.0.0"},
			excluded: []string{"v2", "v1.0.1"},
		},
	}

	for r, tc := range cases {
		t.Run(r, func(t *testing.T) {
			for _, v := range tc.included {
				t.Run(fmt.Sprintf("includes %s", v), func(t *testing.T) {
					assert.True(t, updater.Group{Range: r}.InRange(v))
				})
			}
			for _, v := range tc.excluded {
				t.Run(fmt.Sprintf("excludes %q", v), func(t *testing.T) {
					assert.False(t, updater.Group{Range: r}.InRange(v))
				})
			}
		})
	}
}
