package roll

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmailValidate(t *testing.T) {
	dev := Developer{
		FirstName: "Joe",
		LastName:  "Developer",
		Email:     "foo@dev.com",
	}
	assert.True(t, dev.validateEmail())

	dev.Email = "foo/bar@bar.com"
	assert.False(t, dev.validateEmail())

	dev.Email = ""
	assert.False(t, dev.validateEmail())
}

func TestFirstNameValidate(t *testing.T) {

	dev := Developer{
		FirstName: "Joe",
		LastName:  "Developer",
		Email:     "foo@dev.com",
	}

	assert.True(t, dev.validateFirstName())

	dev.FirstName = "Joe Dev"
	assert.False(t, dev.validateFirstName())

	dev.FirstName = "Joe\"<script>"
	assert.False(t, dev.validateFirstName())
}

func TestLastNameValidate(t *testing.T) {
	dev := Developer{
		FirstName: "Joe",
		LastName:  "Developer",
		Email:     "foo@dev.com",
	}

	assert.True(t, dev.validateLastName())

	dev.LastName = "Van Houten"
	assert.True(t, dev.validateLastName())

	dev.LastName = "v@n hout3n"
	assert.False(t, dev.validateLastName())
}

func TestFullValidation(t *testing.T) {
	dev := Developer{
		FirstName: "Joe Billy Bob",
		LastName:  "D3veloper",
		Email:     "foo/bar@bar.com",
	}

	err := dev.Validate()
	assert.NotNil(t, err)
	msg := err.Error()
	assert.Contains(t, msg, "Email")
	assert.Contains(t, msg, "FirstName")
	assert.Contains(t, msg, "LastName")
}
