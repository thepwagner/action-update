// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package updater_test

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	updater "github.com/thepwagner/action-update/updater"
)

// mockUpdater is an autogenerated mock type for the Updater type
type mockUpdater struct {
	mock.Mock
}

// ApplyUpdate provides a mock function with given fields: _a0, _a1
func (_m *mockUpdater) ApplyUpdate(_a0 context.Context, _a1 updater.Update) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, updater.Update) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Check provides a mock function with given fields: ctx, dep, filter
func (_m *mockUpdater) Check(ctx context.Context, dep updater.Dependency, filter func(string) bool) (*updater.Update, error) {
	ret := _m.Called(ctx, dep, filter)

	var r0 *updater.Update
	if rf, ok := ret.Get(0).(func(context.Context, updater.Dependency, func(string) bool) *updater.Update); ok {
		r0 = rf(ctx, dep, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*updater.Update)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, updater.Dependency, func(string) bool) error); ok {
		r1 = rf(ctx, dep, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Dependencies provides a mock function with given fields: _a0
func (_m *mockUpdater) Dependencies(_a0 context.Context) ([]updater.Dependency, error) {
	ret := _m.Called(_a0)

	var r0 []updater.Dependency
	if rf, ok := ret.Get(0).(func(context.Context) []updater.Dependency); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]updater.Dependency)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
