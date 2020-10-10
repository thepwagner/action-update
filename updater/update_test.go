package updater_test

import (
	"github.com/thepwagner/action-update/updater"
)

const (
	groupName  = "foo"
	baseBranch = "main"
	mockPath   = "github.com/foo/bar"
)

var mockUpdate = updater.Update{
	Path: mockPath,
	Next: "v2.0.0",
}
