package updateaction_test

import (
	"github.com/thepwagner/action-update/actions/updateaction"
	"github.com/thepwagner/action-update/updater"
)

type testEnvironment struct {
	updateaction.Environment
}

func (t *testEnvironment) NewUpdater(string) updater.Updater { return nil }
