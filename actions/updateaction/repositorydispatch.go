package updateaction

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/go-github/v32/github"
	"github.com/thepwagner/action-update/updater"
)

const (
	RepoDispatchActionUpdate = "update-dependency"
)

func (h *handler) RepositoryDispatch(ctx context.Context, evt *github.RepositoryDispatchEvent) error {
	switch evt.GetAction() {
	case RepoDispatchActionUpdate:
		// TODO: is this the right action for this dependency? (e.g. if workflow has multiple action-update-*)
		return h.repoDispatchActionUpdate(ctx, evt)
	default:
		return h.UpdateAll(ctx)
	}
}

func (h *handler) repoDispatchActionUpdate(ctx context.Context, evt *github.RepositoryDispatchEvent) error {
	update, err := unmarshallRepoDispatchUpdate(evt)
	if err != nil {
		return err
	}

	baseBranch := evt.GetRepo().GetDefaultBranch()
	branchName := h.branchNamer.Format(baseBranch, update)
	return h.repoDispatchUpdate(ctx, err, update, baseBranch, branchName)
}

func unmarshallRepoDispatchUpdate(evt *github.RepositoryDispatchEvent) (updater.Update, error) {
	var payload RepoDispatchActionUpdatePayload
	if err := json.Unmarshal(evt.ClientPayload, &payload); err != nil {
		return updater.Update{}, fmt.Errorf("decoding payload: %w", err)
	}
	update := updater.Update{
		Path: payload.Path,
		Next: payload.Next,
	}
	return update, nil
}

func (h *handler) repoDispatchUpdate(ctx context.Context, err error, update updater.Update, baseBranch string, branchName string) error {
	r, err := h.repo()
	if err != nil {
		return fmt.Errorf("getting Repo: %w", err)
	}
	repoUpdater, err := h.repoUpdater(r)
	if err != nil {
		return fmt.Errorf("getting RepoUpdater: %w", err)
	}

	ug := updater.NewUpdateGroup("", update)
	return repoUpdater.Update(ctx, baseBranch, branchName, ug)
}



type RepoDispatchActionUpdatePayload struct {
	Path string `json:"path"`
	Next string `json:"next"`
}
