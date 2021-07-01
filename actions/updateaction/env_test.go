package updateaction_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thepwagner/action-update/actions/updateaction"
	"github.com/thepwagner/action-update/updater"
)

type testEnvironment struct {
	updateaction.Environment
}

func (t *testEnvironment) NewUpdater(string) updater.Updater { return nil }

func TestEnvironment_Ignored(t *testing.T) {
	cases := []struct {
		ignore     string
		ignored    []string
		notIgnored []string
	}{
		{
			ignore:     "foo/*",
			ignored:    []string{"foo/bar"},
			notIgnored: []string{"foo", "foo/bar/bar", "bar/foo", "bar/foo/bar"},
		},
		{
			ignore:     "foo/*\nbar/**",
			ignored:    []string{"bar", "foo/bar", "bar/foo", "bar/foo/bar"},
			notIgnored: []string{"foo/bar/bar"},
		},
		{
			ignore:     "foo/**/*.bar",
			ignored:    []string{"foo/a.bar", "foo/bar/b.bar"},
			notIgnored: []string{"a.bar", "foo/bar", "bar/foo", "bar/foo/bar", "bar", "foo/bar/bar"},
		},
		{
			ignore:     "foo/**/*.bar",
			ignored:    []string{"foo/a.bar", "foo/bar/b.bar"},
			notIgnored: []string{"a.bar", "foo/bar", "bar/foo", "bar/foo/bar", "bar", "foo/bar/bar"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.ignore, func(t *testing.T) {
			e := updateaction.Environment{InputIgnore: tc.ignore}
			for _, i := range tc.ignored {
				assert.True(t, e.Ignored(i))
			}
			for _, i := range tc.notIgnored {
				assert.False(t, e.Ignored(i))
			}
		})
	}
}
