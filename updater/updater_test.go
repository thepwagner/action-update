package updater_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/action-update/updater"
)

//go:generate mockery --outpkg updater_test --output . --testonly --name Updater --structname mockUpdater --filename mockupdater_test.go
//go:generate mockery --outpkg updater_test --output . --testonly --name Repo --structname mockRepo --filename mockrepo_test.go

func TestRepoUpdater_Update(t *testing.T) {
	r := &mockRepo{}
	u := &mockUpdater{}
	ru := updater.NewRepoUpdater(r, u)
	ctx := context.Background()

	branch := setupMockUpdate(ctx, r, u, mockUpdate)

	err := ru.Update(ctx, baseBranch, branch, updater.NewUpdateGroup("", mockUpdate))
	require.NoError(t, err)
	r.AssertExpectations(t)
	u.AssertExpectations(t)
}

func setupMockUpdate(ctx context.Context, r *mockRepo, u *mockUpdater, up updater.Update) string {
	branch := fmt.Sprintf("action-update-go/main/%s/%s", up.Path, up.Next)
	r.On("NewBranch", baseBranch, branch).Return(nil)
	u.On("ApplyUpdate", ctx, up).Return(nil)
	r.On("Push", ctx, updater.NewUpdateGroup("", up)).Return(nil)
	return branch
}

func TestRepoUpdater_UpdateAll_NoChanges(t *testing.T) {
	r := &mockRepo{}
	u := &mockUpdater{}
	ru := updater.NewRepoUpdater(r, u)
	ctx := context.Background()

	r.On("SetBranch", baseBranch).Return(nil)
	dep := updater.Dependency{Path: mockUpdate.Path, Version: mockUpdate.Previous}
	u.On("Dependencies", ctx).Return([]updater.Dependency{dep}, nil)
	r.On("ExistingUpdates", ctx, baseBranch).Return(nil, nil)
	u.On("Check", ctx, dep, mock.Anything).Return(nil, nil)

	err := ru.UpdateAll(ctx, baseBranch)
	require.NoError(t, err)
	r.AssertExpectations(t)
	u.AssertExpectations(t)
}

func TestRepoUpdater_UpdateAll_Update(t *testing.T) {
	r := &mockRepo{}
	u := &mockUpdater{}
	ru := updater.NewRepoUpdater(r, u)
	ctx := context.Background()

	r.On("SetBranch", baseBranch).Return(nil)
	dep := updater.Dependency{Path: mockUpdate.Path, Version: mockUpdate.Previous}
	u.On("Dependencies", ctx).Return([]updater.Dependency{dep}, nil)
	r.On("ExistingUpdates", ctx, baseBranch).Return(nil, nil)
	availableUpdate := mockUpdate // avoid pointer to shared reference
	u.On("Check", ctx, dep, mock.Anything).Return(&availableUpdate, nil)
	setupMockUpdate(ctx, r, u, mockUpdate) // delegates to .Update()

	err := ru.UpdateAll(ctx, baseBranch)
	require.NoError(t, err)
	r.AssertExpectations(t)
	u.AssertExpectations(t)
}

func TestRepoUpdater_UpdateAll_Multiple(t *testing.T) {
	r := &mockRepo{}
	u := &mockUpdater{}
	ru := updater.NewRepoUpdater(r, u)
	ctx := context.Background()

	r.On("SetBranch", baseBranch).Return(nil)
	dep := updater.Dependency{Path: mockUpdate.Path, Version: mockUpdate.Previous}
	otherDep := updater.Dependency{Path: "github.com/foo/baz", Version: mockUpdate.Previous}
	u.On("Dependencies", ctx).Return([]updater.Dependency{dep, otherDep}, nil)
	r.On("ExistingUpdates", ctx, baseBranch).Return(nil, nil)
	availableUpdate := mockUpdate // avoid pointer to shared reference
	u.On("Check", ctx, dep, mock.Anything).Return(&availableUpdate, nil)
	otherUpdate := updater.Update{Path: "github.com/foo/baz", Next: "v3.0.0"}
	u.On("Check", ctx, otherDep, mock.Anything).Return(&otherUpdate, nil)
	setupMockUpdate(ctx, r, u, mockUpdate)  // delegates to .Update()
	setupMockUpdate(ctx, r, u, otherUpdate) // delegates to .Update()

	err := ru.UpdateAll(ctx, baseBranch)
	require.NoError(t, err)
	r.AssertExpectations(t)
	u.AssertExpectations(t)
}

func TestRepoUpdater_UpdateAll_MultipleWithError(t *testing.T) {
	r := &mockRepo{}
	u := &mockUpdater{}
	ru := updater.NewRepoUpdater(r, u)
	ctx := context.Background()

	r.On("SetBranch", baseBranch).Return(nil)
	dep := updater.Dependency{Path: mockUpdate.Path, Version: mockUpdate.Previous}
	otherDep := updater.Dependency{Path: "github.com/foo/baz", Version: mockUpdate.Previous}
	u.On("Dependencies", ctx).Return([]updater.Dependency{dep, otherDep}, nil)
	r.On("ExistingUpdates", ctx, baseBranch).Return(nil, nil)
	availableUpdate := mockUpdate // avoid pointer to shared reference
	u.On("Check", ctx, dep, mock.Anything).Return(&availableUpdate, nil)
	otherUpdate := updater.Update{Path: "github.com/foo/baz", Next: "v3.0.0"}
	u.On("Check", ctx, otherDep, mock.Anything).Return(&otherUpdate, nil)
	branch := fmt.Sprintf("action-update-go/main/%s/%s", mockUpdate.Path, mockUpdate.Next)
	r.On("NewBranch", baseBranch, branch).Return(nil)
	u.On("ApplyUpdate", ctx, mockUpdate).Return(fmt.Errorf("kaboom"))
	setupMockUpdate(ctx, r, u, otherUpdate) // delegates to .Update()

	err := ru.UpdateAll(ctx, baseBranch)
	require.NoError(t, err)
	r.AssertExpectations(t)
	u.AssertExpectations(t)
}

func TestRepoUpdater_UpdateAll_MultipleGrouped(t *testing.T) {
	group := &updater.Group{Name: groupName, Pattern: "github.com/foo"}
	err := group.Validate()
	require.NoError(t, err)

	r := &mockRepo{}
	u := &mockUpdater{}
	ru := updater.NewRepoUpdater(r, u, updater.WithGroups(group))
	ctx := context.Background()

	r.On("SetBranch", baseBranch).Return(nil)
	dep := updater.Dependency{Path: mockUpdate.Path, Version: mockUpdate.Previous}
	otherDep := updater.Dependency{Path: "github.com/foo/baz", Version: mockUpdate.Previous}
	u.On("Dependencies", ctx).Return([]updater.Dependency{dep, otherDep}, nil)
	r.On("ExistingUpdates", ctx, baseBranch).Return(nil, nil)
	availableUpdate := mockUpdate // avoid pointer to shared reference
	u.On("Check", ctx, dep, mock.Anything).Return(&availableUpdate, nil)
	otherUpdate := updater.Update{Path: "github.com/foo/baz", Next: "v3.0.0"}
	u.On("Check", ctx, otherDep, mock.Anything).Return(&otherUpdate, nil)

	r.On("NewBranch", baseBranch, "action-update-go/main/foo").Times(1).Return(nil)
	u.On("ApplyUpdate", ctx, mock.Anything).Times(2).Return(nil)
	r.On("Push", ctx, mock.Anything, mock.Anything).Times(1).Return(nil)

	err = ru.UpdateAll(ctx, baseBranch)
	require.NoError(t, err)
	r.AssertExpectations(t)
	u.AssertExpectations(t)
}

func TestRepoUpdater_UpdateAll_Scripts(t *testing.T) {
	cases := []*updater.Group{
		{
			Name:      groupName,
			Pattern:   "github.com/foo",
			PreScript: `echo "sup" && touch token`,
		},
		{
			Name:       groupName,
			Pattern:    "github.com/foo",
			PostScript: `echo "sup" && touch token`,
		},
	}

	for _, group := range cases {
		err := group.Validate()
		require.NoError(t, err)

		tmpDir := t.TempDir()
		tokenPath := filepath.Join(tmpDir, "token")
		r := &mockRepo{}
		u := &mockUpdater{}
		ru := updater.NewRepoUpdater(r, u, updater.WithGroups(group))
		ctx := context.Background()

		r.On("SetBranch", baseBranch).Return(nil)
		dep := updater.Dependency{Path: mockUpdate.Path, Version: mockUpdate.Previous}
		u.On("Dependencies", ctx).Return([]updater.Dependency{dep}, nil)
		r.On("ExistingUpdates", ctx, baseBranch).Return(nil, nil)
		availableUpdate := mockUpdate // avoid pointer to shared reference
		u.On("Check", ctx, dep, mock.Anything).Return(&availableUpdate, nil)
		r.On("NewBranch", baseBranch, "action-update-go/main/foo").Times(1).Return(nil)
		u.On("ApplyUpdate", ctx, mock.Anything).Times(1).Return(nil)
		r.On("Push", ctx, mock.Anything, mock.Anything).Times(1).Return(nil)
		r.On("Root").Return(tmpDir)

		err = ru.UpdateAll(ctx, baseBranch)
		require.NoError(t, err)
		r.AssertExpectations(t)
		u.AssertExpectations(t)
		_, err = os.Stat(tokenPath)
		require.NoError(t, err)
	}
}

func TestRepoUpdater_UpdateAll_CoolDown(t *testing.T) {
	// Group with 1 day cooldown
	group := &updater.Group{Name: groupName, Pattern: "github.com/foo", CoolDown: "1D"}
	err := group.Validate()
	require.NoError(t, err)

	r := &mockRepo{}
	u := &mockUpdater{}
	ru := updater.NewRepoUpdater(r, u, updater.WithGroups(group))
	ctx := context.Background()
	r.On("SetBranch", baseBranch).Return(nil)
	dep := updater.Dependency{Path: mockUpdate.Path, Version: mockUpdate.Previous}
	otherDep := updater.Dependency{Path: "github.com/foo/baz", Version: mockUpdate.Previous}
	u.On("Dependencies", ctx).Return([]updater.Dependency{dep, otherDep}, nil)

	// Existing update within 30m means update is not checked:
	existingUpdates := updater.ExistingUpdates{
		{Group: updater.UpdateGroup{Name: groupName}, Merged: true, LastUpdate: time.Now().Add(-1 * time.Hour)},
	}
	r.On("ExistingUpdates", ctx, baseBranch).Return(existingUpdates, nil)

	err = ru.UpdateAll(ctx, baseBranch)
	require.NoError(t, err)
	r.AssertExpectations(t)
	u.AssertExpectations(t)
}
