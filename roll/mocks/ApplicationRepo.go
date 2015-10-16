package mocks

import "github.com/xtraclabs/roll/roll"
import "github.com/stretchr/testify/mock"

type ApplicationRepo struct {
	mock.Mock
}

func (_m *ApplicationRepo) CreateApplication(app *roll.Application) error {
	ret := _m.Called(app)

	var r0 error
	if rf, ok := ret.Get(0).(func(*roll.Application) error); ok {
		r0 = rf(app)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ApplicationRepo) UpdateApplication(app *roll.Application) error {
	ret := _m.Called(app)

	var r0 error
	if rf, ok := ret.Get(0).(func(*roll.Application) error); ok {
		r0 = rf(app)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ApplicationRepo) RetrieveApplication(clientID string) (*roll.Application, error) {
	ret := _m.Called(clientID)

	var r0 *roll.Application
	if rf, ok := ret.Get(0).(func(string) *roll.Application); ok {
		r0 = rf(clientID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*roll.Application)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(clientID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *ApplicationRepo) ListApplications() ([]roll.Application, error) {
	ret := _m.Called()

	var r0 []roll.Application
	if rf, ok := ret.Get(0).(func() []roll.Application); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]roll.Application)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
