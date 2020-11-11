package updateaction

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v32/github"
	"github.com/sirupsen/logrus"
	"github.com/thepwagner/action-update/repo"
)

func (h *handler) Release(ctx context.Context, evt *github.ReleaseEvent) error {
	if evt.GetAction() != "released" {
		logrus.WithField("action", evt.GetAction()).Info("ignoring release event")
		return nil
	}
	notifyRepos := h.cfg.ReleaseDispatchRepos()
	logrus.WithField("repos", len(notifyRepos)).Info("notifying repositories of release")
	if len(notifyRepos) == 0 {
		return nil
	}

	dispatchOpts, err := releaseDispatchOptions(evt)
	if err != nil {
		return err
	}

	gh := repo.NewGitHubClient(h.cfg.GitHubToken)
	for _, notifyRepo := range notifyRepos {
		notifyRepoParts := strings.SplitN(notifyRepo, "/", 2)
		owner := notifyRepoParts[0]
		name := notifyRepoParts[1]
		if _, _, err := gh.Repositories.Dispatch(ctx, owner, name, dispatchOpts); err != nil {
			logrus.WithError(err).Warn("error dispatching update")
		}
	}

	return nil
}

func releaseDispatchOptions(evt *github.ReleaseEvent) (github.DispatchRequestOptions, error) {
	payload, err := json.Marshal(&RepoDispatchActionUpdatePayload{
		Path: fmt.Sprintf("github.com/%s", evt.GetRepo().GetFullName()),
		Next: evt.GetRelease().GetTagName(),
	})
	if err != nil {
		return github.DispatchRequestOptions{}, fmt.Errorf("serializing payload: %w", err)
	}
	clientPayload := json.RawMessage(payload)
	return github.DispatchRequestOptions{
		EventType:     RepoDispatchActionUpdate,
		ClientPayload: &clientPayload,
	}, nil
}
