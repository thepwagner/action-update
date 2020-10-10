package repo_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/action-update/repo"
)

func TestNewGitHubRepo(t *testing.T) {
	gr := initGitRepo(t, plumbing.NewBranchReferenceName(branchName))

	gh, err := repo.NewGitHubRepo(gr, testKey, "foo/bar", "")
	require.NoError(t, err)
	assert.NotNil(t, gh)
}

func TestGitHubRepo_ExistingUpdates(t *testing.T) {
	gr := initGitRepo(t, plumbing.NewBranchReferenceName(branchName))

	gh, err := repo.NewGitHubRepo(gr, []byte(""), "thepwagner/action-update", "")
	require.NoError(t, err)

	existing, err := gh.ExistingUpdates(context.Background(), "main")
	require.NoError(t, err)
	assert.NotEmpty(t, existing)
}

func tokenOrSkip(t *testing.T) string {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		t.Skip("set GITHUB_TOKEN")
	}
	return token
}
