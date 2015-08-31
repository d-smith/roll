package roll

import (
	"testing"
)

func TestEmailValidate(t *testing.T) {
	if !ValidateEmail("foo@dev.com") {
		t.Fail()
	}

	if ValidateEmail("foo/bar@bar.com") {
		t.Fail()
	}
}
