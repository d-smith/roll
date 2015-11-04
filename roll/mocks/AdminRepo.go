package mocks

import "github.com/stretchr/testify/mock"

type AdminRepo struct {
	mock.Mock
}

func (_m *AdminRepo) IsAdmin(subject string) (bool, error) {
	ret := _m.Called(subject)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(subject)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(subject)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
