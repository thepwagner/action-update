package updater_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/action-update/updater"
)

func TestGroup_CoolDownDuration(t *testing.T) {
	g := updater.Group{Name: "test", Pattern: "test"}

	cases := map[string]time.Duration{
		"P1D": 24 * time.Hour,
		"1D":  24 * time.Hour,
		"1W":  7 * 24 * time.Hour,
	}

	for in, expected := range cases {
		t.Run(in, func(t *testing.T) {
			g.CoolDown = in
			err := g.Validate()
			require.NoError(t, err)
			assert.Equal(t, expected, g.CoolDownDuration())
		})
	}
}

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
		"< 1.0.0": {
			included: []string{"v0.1"},
			excluded: []string{"v2", "v1.0.1", "v1", "v1.0.0"},
		},
		"<= v1": {
			included: []string{"v0.1", "v1", "v1.0.0"},
			excluded: []string{"v2", "v1.0.1"},
		},
		">=v2 , < v3": {
			included: []string{"v2", "v2.0.0"},
			excluded: []string{"v0.1", "v1", "v1.0.0", "v1.0.1", "v3"},
		},
		"": {
			included: []string{"v0.1", "v1", "v1.0.0", "v2", "v1.0.1"},
		},
	}

	for r, tc := range cases {
		t.Run(r, func(t *testing.T) {
			u := &updater.Group{Range: r}
			for _, v := range tc.included {
				t.Run(fmt.Sprintf("includes %s", v), func(t *testing.T) {
					assert.True(t, u.InRange(v))
				})
			}
			for _, v := range tc.excluded {
				t.Run(fmt.Sprintf("excludes %q", v), func(t *testing.T) {
					assert.False(t, u.InRange(v))
				})
			}
		})
	}
}
