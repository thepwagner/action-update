package actions_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-github/v32/github"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/action-update/actions"
)

func TestEnvironment_LogLevel(t *testing.T) {
	cases := map[string]logrus.Level{
		"":        logrus.InfoLevel,
		"invalid": logrus.InfoLevel,
		"warn":    logrus.WarnLevel,
	}

	for in, lvl := range cases {
		t.Run(fmt.Sprintf("parse %q", in), func(t *testing.T) {
			e := actions.Environment{InputLogLevel: in}
			assert.Equal(t, lvl, e.LogLevel())
		})
	}
}

func TestEnvironment_ParseEvent_Noop(t *testing.T) {
	e := actions.Environment{GitHubEventName: "schedule"}
	evt, err := e.ParseEvent()
	require.NoError(t, err)
	assert.Nil(t, evt)
}

func TestEnvironment_ParseEvent(t *testing.T) {
	e := actions.Environment{
		GitHubEventName: "issue_comment",
		GitHubEventPath: testIssueComment(t, "test"),
	}

	evt, err := e.ParseEvent()
	require.NoError(t, err)

	ic, ok := evt.(*github.IssueCommentEvent)
	if assert.True(t, ok) {
		assert.Equal(t, "test", ic.GetComment().GetBody())
	}
}

func testIssueComment(t *testing.T, body string) string {
	eventPath := filepath.Join(t.TempDir(), "event.json")
	err := ioutil.WriteFile(eventPath, []byte(fmt.Sprintf(`{"comment":{"body":%q}}`, body)), 0600)
	require.NoError(t, err)
	return eventPath
}

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
			ignored:    []string{"foo/bar", "bar/foo", "bar/foo/bar"},
			notIgnored: []string{"bar", "foo/bar/bar"},
		},
		{
			ignore:     "foo/**/*.bar",
			ignored:    []string{"foo/a.bar", "foo/bar/b.bar"},
			notIgnored: []string{"a.bar", "foo/bar", "bar/foo", "bar/foo/bar", "bar", "foo/bar/bar"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.ignore, func(t *testing.T) {
			e := actions.Environment{InputIgnore: tc.ignore}
			for _, i := range tc.ignored {
				assert.True(t, e.Ignored(i))
			}
			for _, i := range tc.notIgnored {
				assert.False(t, e.Ignored(i))
			}
		})
	}
}
