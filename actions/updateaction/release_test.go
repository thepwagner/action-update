package updateaction_test

import (
	"context"
	"testing"

	"github.com/google/go-github/v37/github"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/action-update/actions/updateaction"
)

func TestHandler_Release_WrongAction(t *testing.T) {
	handlers := updateaction.NewHandlers(&testEnvironment{})
	err := handlers.Release(context.Background(), &github.ReleaseEvent{
		Action: github.String("prereleased"),
	})
	require.NoError(t, err)
}

func TestHandler_Release_NoRepos(t *testing.T) {
	handlers := updateaction.NewHandlers(&testEnvironment{})
	err := handlers.Release(context.Background(), &github.ReleaseEvent{
		Action: github.String("released"),
	})
	require.NoError(t, err)
}
