package mocks

import "github.com/xtraclabs/roll/roll"
import "github.com/stretchr/testify/mock"

type DeveloperRepo struct {
	mock.Mock
}

func (_m *DeveloperRepo) RetrieveDeveloper(email string) (*roll.Developer, error) {
	ret := _m.Called(email)

	var r0 *roll.Developer
	if rf, ok := ret.Get(0).(func(string) *roll.Developer); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*roll.Developer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *DeveloperRepo) StoreDeveloper(_a0 *roll.Developer) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*roll.Developer) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
