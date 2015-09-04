package mocks

import "github.com/stretchr/testify/mock"

type SecretsRepo struct {
	mock.Mock
}

func (_m *SecretsRepo) StoreKeysForApp(appkey string, privateKey string, publicKey string) error {
	ret := _m.Called(appkey, privateKey, publicKey)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(appkey, privateKey, publicKey)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *SecretsRepo) RetrievePrivateKeyForApp(appkey string) (string, error) {
	ret := _m.Called(appkey)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(appkey)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(appkey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *SecretsRepo) RetrievePublicKeyForApp(appkey string) (string, error) {
	ret := _m.Called(appkey)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(appkey)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(appkey)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
