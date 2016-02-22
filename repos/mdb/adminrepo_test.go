// +build integration

package mdb

import (
	"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("ROLL_DBADDRESS", "127.0.0.1:3306")
	os.Setenv("ROLL_DBPASSWORD", "rollpw")
	os.Setenv("ROLL_DBUSER", "rolluser")

	flag.Parse()
	os.Exit(m.Run())
}

func TestIsAdmin(t *testing.T) {
	adminRepo := NewMBDAdminRepo()
	admin, err := adminRepo.IsAdmin("foobar")
	if assert.Nil(t, err) {
		assert.False(t, admin)
	}
}
