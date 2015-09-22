package roll

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateAppName(t *testing.T) {
	var app = Application{
		ApplicationName: "Most excellent app",
		DeveloperEmail:  "jane@someplace.com",
		RedirectURI:     "http://google.com/login_callback",
		LoginProvider:   "xtrac://bigiron:9000",
	}

	assert.True(t, app.validateApplicationName())

	app.ApplicationName = "<script/>"
	assert.False(t, app.validateApplicationName())

	app.ApplicationName = "Most excellent app v1.1"
	assert.True(t, app.validateApplicationName())

}

func TestValidateDeveloperEmail(t *testing.T) {
	var app = Application{
		ApplicationName: "Most excellent app",
		DeveloperEmail:  "jane@someplace.com",
		RedirectURI:     "http://google.com/login_callback",
		LoginProvider:   "xtrac://bigiron:9000",
	}

	assert.True(t, app.validateDeveloperEmail())

	app.DeveloperEmail = "foo@<script/>"
	assert.False(t, app.validateDeveloperEmail())

	app.DeveloperEmail = "foo.com"
	assert.False(t, app.validateDeveloperEmail())

	app.DeveloperEmail = "@foo.com"
	assert.False(t, app.validateDeveloperEmail())
}

func TestValidateLoginProvider(t *testing.T) {
	var app = Application{
		ApplicationName: "Most excellent app",
		DeveloperEmail:  "jane@someplace.com",
		RedirectURI:     "http://google.com/login_callback",
		LoginProvider:   "xtrac://bigiron:9000",
	}

	assert.True(t, app.validateLoginProvider())

	app.LoginProvider = "foo://localhost:9000"
	assert.False(t, app.validateLoginProvider())
}

func TestValidateRedirectURI(t *testing.T) {
	var app = Application{
		ApplicationName: "Most excellent app",
		DeveloperEmail:  "jane@someplace.com",
		RedirectURI:     "http://google.com/login_callback",
		LoginProvider:   "xtrac://bigiron:9000",
	}

	assert.True(t, app.validateRedirectURI())

	app.RedirectURI = "huh?"
	assert.False(t, app.validateRedirectURI())

	app.RedirectURI = "https://google.com/"
	assert.True(t, app.validateRedirectURI())

}

func TestApplicationFullValidationOK(t *testing.T) {
	var app = Application{
		ApplicationName: "Most excellent app",
		DeveloperEmail:  "jane@someplace.com",
		RedirectURI:     "http://google.com/login_callback",
		LoginProvider:   "xtrac://bigiron:9000",
	}

	err := app.Validate()
	assert.Nil(t, err)
}

func TestApplicationFullValidationAllBad(t *testing.T) {
	var app = Application{
		ApplicationName: "Most excelle<script/>nt app",
		DeveloperEmail:  "jane@some<script/>place.com",
		RedirectURI:     "xxx",
		LoginProvider:   "xxx",
	}

	err := app.Validate()
	assert.NotNil(t, err)
	msg := err.Error()
	println(msg)
	assert.Contains(t, msg, "ApplicationName")
	assert.Contains(t, msg, "DeveloperEmail")
	assert.Contains(t, msg, "LoginProvider")
	assert.Contains(t, msg, "RedirectURI")
}
