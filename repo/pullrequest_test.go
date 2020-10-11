package repo_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/action-update/repo"
	"github.com/thepwagner/action-update/updater"
)

var (
	awsSdkGo13417 = updater.Update{
		Path:     "github.com/aws/aws-sdk-go",
		Previous: "v1.34.16",
		Next:     "v1.34.17",
	}
	fooBar987 = updater.Update{
		Path:     "github.com/foo/bar",
		Previous: "v0.4.1",
		Next:     "v99.88.77",
	}
	testKey = []byte{1, 2, 3, 4}
)

func TestGitHubPullRequestContent_Generate(t *testing.T) {
	token := tokenOrSkip(t)
	client := repo.NewGitHubClient(token)
	gen := repo.NewGitHubPullRequestContent(client, testKey)

	title, body, err := gen.Generate(context.Background(), updater.NewUpdateGroup("", awsSdkGo13417))
	require.NoError(t, err)
	assert.Equal(t, "Update github.com/aws/aws-sdk-go from v1.34.16 to v1.34.17", title)
	assert.Equal(t, strings.TrimSpace(`
Here is github.com/aws/aws-sdk-go v1.34.17, I hope it works.

[changelog](https://github.com/aws/aws-sdk-go/blob/v1.34.17/CHANGELOG.md)

<!--::action-update-go::
{"signed":{"updates":[{"path":"github.com/aws/aws-sdk-go","previous":"v1.34.16","next":"v1.34.17"}]},"signature":"0nxLHGFk/K3Iyi31ArR6wfS9nDMxSvnjcf4i4AeYhj3LWmiwdMDMySMLAkZ1nM/zuVWsENE3zfHy8cC6/6akGg=="}
-->
`), strings.TrimSpace(body))
}

func TestGitHubPullRequestContent_ParseBody(t *testing.T) {
	token := tokenOrSkip(t)
	client := repo.NewGitHubClient(token)
	gen := repo.NewGitHubPullRequestContent(client, testKey)

	body := `
<!--::action-update-go::
{"signed":{"updates":[{"path":"github.com/aws/aws-sdk-go","previous":"v1.34.16","next":"v1.34.17"}]},"signature":"0nxLHGFk/K3Iyi31ArR6wfS9nDMxSvnjcf4i4AeYhj3LWmiwdMDMySMLAkZ1nM/zuVWsENE3zfHy8cC6/6akGg=="}
-->`
	parsed := gen.ParseBody(body)
	assert.Equal(t, []updater.Update{awsSdkGo13417}, parsed.Updates)
	assert.Equal(t, "", parsed.Name)
}

func TestGitHubPullRequestContent_GenerateNoChangeLog(t *testing.T) {
	token := tokenOrSkip(t)
	client := repo.NewGitHubClient(token)
	gen := repo.NewGitHubPullRequestContent(client, testKey)

	title, body, err := gen.Generate(context.Background(), updater.NewUpdateGroup("", fooBar987))
	require.NoError(t, err)
	assert.Equal(t, "Update github.com/foo/bar from v0.4.1 to v99.88.77", title)
	assert.Equal(t, strings.TrimSpace(`
Here is github.com/foo/bar v99.88.77, I hope it works.

<!--::action-update-go::
{"signed":{"updates":[{"path":"github.com/foo/bar","previous":"v0.4.1","next":"v99.88.77"}]},"signature":"hSLvci96ReSaNrsSJ/yw9IsK9AfAvvXHWtJDlTh8TtZZth2vfT7/66BPmKDGb2GYQDNvDavFLOgtkHeWWT5ZTg=="}
-->
`), strings.TrimSpace(body))
}

func TestGitHubPullRequestContent_GenerateMultiple(t *testing.T) {
	token := tokenOrSkip(t)
	client := repo.NewGitHubClient(token)
	gen := repo.NewGitHubPullRequestContent(client, testKey)

	ug := updater.NewUpdateGroup("my-awesome-group", awsSdkGo13417, fooBar987)
	title, body, err := gen.Generate(context.Background(), ug)
	require.NoError(t, err)
	assert.Equal(t, "Dependency Updates", title)
	assert.Equal(t, strings.TrimSpace(`
Here are some updates, I hope they work.

#### github.com/aws/aws-sdk-go@v1.34.17

[changelog](https://github.com/aws/aws-sdk-go/blob/v1.34.17/CHANGELOG.md)

#### github.com/foo/bar@v99.88.77

<!--::action-update-go::
{"signed":{"name":"my-awesome-group","updates":[{"path":"github.com/aws/aws-sdk-go","previous":"v1.34.16","next":"v1.34.17"},{"path":"github.com/foo/bar","previous":"v0.4.1","next":"v99.88.77"}]},"signature":"kpmZqS8mPeaKpl3T9cLsHGutsSaG3ZkXA15gOwAg6H5kBadk5H496Zev+CInIuWuK6EyyPpxQHr9MHthoC8Xdw=="}
-->
`), strings.TrimSpace(body))
}
