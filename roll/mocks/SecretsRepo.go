package mocks

import "github.com/stretchr/testify/mock"

type SecretsRepo struct {
	mock.Mock
}

func (_m *SecretsRepo) StoreKeysForApp(appid string, privateKey string, publicKey string) error {
	ret := _m.Called(appid, privateKey, publicKey)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(appid, privateKey, publicKey)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
