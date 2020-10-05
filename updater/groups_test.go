package updater_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/action-update/updater"
)

func TestParseGroups(t *testing.T) {
	cases := map[string]struct {
		in       string
		err      string
		expected updater.Groups
	}{
		"empty": {
			in:       ``,
			expected: updater.Groups{},
		},
		"single": {
			in: `---
- name: foo
  pattern: github.com/thepwagner
  frequency: weekly
  range: ">=v1.4.0, <v2"`,
			expected: updater.Groups{{
				Name:      "foo",
				Pattern:   "github.com/thepwagner",
				Frequency: "weekly",
				Range:     ">=v1.4.0, <v2",
			}},
		},
		"multiple": {
			in: `---
- name: foo
  pattern: github.com
- name: bar
  pattern: gopkg.included`,
			expected: updater.Groups{
				{
					Name:    "foo",
					Pattern: "github.com",
				},
				{
					Name:    "bar",
					Pattern: "gopkg.included",
				},
			},
		},
		"bad yaml": {
			in:  `{"worst_yaml":"is_json"}`,
			err: "yaml: unmarshal errors:\n  line 1: cannot unmarshal !!map into updater.Groups",
		},
		"duplicate name": {
			in: `---
- name: foo
  pattern: foo
- name: foo
  pattern: foo`,
			err: `duplicate group name: "foo"`,
		},
		"regexp pattern": {
			in: `---
- name: foo
  pattern: /.*pwagner.*/`,
			expected: updater.Groups{{
				Name:    "foo",
				Pattern: "/.*pwagner.*/",
			}},
		},
	}

	for label, tc := range cases {
		t.Run(label, func(t *testing.T) {
			groups, err := updater.ParseGroups(tc.in)
			if tc.err != "" {
				assert.EqualError(t, err, tc.err)
				assert.Nil(t, groups)
			} else {
				require.NoError(t, err)
				if assert.Equal(t, len(tc.expected), len(groups)) {
					for i, g := range groups {
						assert.Equal(t, tc.expected[i].Name, g.Name)
						assert.Equal(t, tc.expected[i].Pattern, g.Pattern)
						assert.Equal(t, tc.expected[i].Frequency, g.Frequency)
						assert.Equal(t, tc.expected[i].Range, g.Range)
					}
				}
			}
		})
	}
}

func TestGroups_GroupDependencies(t *testing.T) {
	groups := updater.Groups{
		{
			Name:    "contains foo",
			Pattern: "/.*foo.*/",
		},
		{
			Name:    "prefix bar",
			Pattern: "bar",
		},
	}
	err := groups.Validate()
	require.NoError(t, err)

	byGroupName, ungrouped := groups.GroupDependencies([]updater.Dependency{
		{Path: "foo at the start"},
		{Path: "included the foo middle"},
		{Path: "final foo"},
		{Path: "bar"},
		{Path: "bar at the start"},
		{Path: "final bar"},
		{Path: "no match"},
		{Path: "emoji 😀"},
	})

	assert.Equal(t, map[string][]updater.Dependency{
		"contains foo": {
			{Path: "foo at the start"},
			{Path: "included the foo middle"},
			{Path: "final foo"},
		},
		"prefix bar": {
			{Path: "bar"},
			{Path: "bar at the start"},
		},
	}, byGroupName)
	assert.Equal(t, []updater.Dependency{
		{Path: "final bar"},
		{Path: "no match"},
		{Path: "emoji 😀"},
	}, ungrouped)
}
