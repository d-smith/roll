package mocks

import "github.com/xtraclabs/roll/roll"
import "github.com/stretchr/testify/mock"

type ApplicationRepo struct {
	mock.Mock
}

func (_m *ApplicationRepo) StoreApplication(app *roll.Application) error {
	ret := _m.Called(app)

	var r0 error
	if rf, ok := ret.Get(0).(func(*roll.Application) error); ok {
		r0 = rf(app)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *ApplicationRepo) RetrieveApplication(apiKey string) (*roll.Application, error) {
	ret := _m.Called(apiKey)

	var r0 *roll.Application
	if rf, ok := ret.Get(0).(func(string) *roll.Application); ok {
		r0 = rf(apiKey)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*roll.Application)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(apiKey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
