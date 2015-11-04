package http

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/xtraclabs/roll/roll/mocks"
	"testing"
)

func TestAdminScopeWithNoAdminPrivileges(t *testing.T) {
	core, coreConfig := NewTestCore()

	appRepoMock := coreConfig.AdminRepo.(*mocks.AdminRepo)
	appRepoMock.On("IsAdmin", "foobar").Return(false, nil)

	granted, err := grantAdminScope(core, "foobar")
	assert.Nil(t, err)
	assert.False(t, granted)

}

func TestAdminScopeWithAdminPrivileges(t *testing.T) {
	core, coreConfig := NewTestCore()

	appRepoMock := coreConfig.AdminRepo.(*mocks.AdminRepo)
	appRepoMock.On("IsAdmin", "foobar").Return(true, nil)

	granted, err := grantAdminScope(core, "foobar")
	assert.Nil(t, err)
	assert.True(t, granted)

}

func TestAdminScopeWithRepoError(t *testing.T) {
	core, coreConfig := NewTestCore()

	appRepoMock := coreConfig.AdminRepo.(*mocks.AdminRepo)
	appRepoMock.On("IsAdmin", "foobar").Return(true, errors.New("boom boom out go the lights"))

	_, err := grantAdminScope(core, "foobar")
	assert.NotNil(t, err)
}
