package roll

import (
	"bytes"
	"errors"
	"github.com/xtraclabs/roll/login"
	"net/url"
	"regexp"
	"strings"
)

//Application represents the data associated with an application that is exposed via the REST API
type Application struct {
	DeveloperEmail   string `json:"developerEmail"`
	DeveloperID      string `json:developerID`
	ClientID         string `json:"clientID"`
	ApplicationName  string `json:"applicationName"`
	ClientSecret     string `json:"clientSecret"`
	RedirectURI      string `json:"redirectURI"`
	LoginProvider    string `json:"loginProvider"`
	JWTFlowPublicKey string `json:"jwtFlowPublicKey"`
}

var appName = regexp.MustCompile(`^([a-zA-Z'-.0-9]\s*)+$`)

func (a *Application) validateDeveloperEmail() bool {
	return validEmail.MatchString(a.DeveloperEmail)
}

func (a *Application) validateApplicationName() bool {
	return appName.MatchString(a.ApplicationName)
}

func (a *Application) validateRedirectURI() bool {
	parsed, err := url.Parse(a.RedirectURI)
	if err != nil {
		return false
	}

	if parsed.Scheme == "" || parsed.Host == "" || parsed.Path == "" {
		return false
	}

	if !strings.HasPrefix(parsed.Scheme, "http") {
		return false
	}

	return true
}

func (a *Application) validateLoginProvider() bool {
	parsed, err := url.Parse(a.LoginProvider)
	if err != nil {
		return false
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return false
	}

	return login.SupportedProvider(parsed.Scheme)
}

func (a *Application) Validate() error {
	var valid = true
	var err error

	bs := bytes.NewBufferString("Fields with invalid content: ")

	if !a.validateApplicationName() {
		valid = false
		bs.WriteString("ApplicationName ")
	}

	if !a.validateDeveloperEmail() {
		valid = false
		bs.WriteString("DeveloperEmail ")
	}

	if !a.validateRedirectURI() {
		valid = false
		bs.WriteString("RedirectURI ")
	}

	if !a.validateLoginProvider() {
		valid = false
		bs.WriteString("LoginProvider ")
	}

	if !valid {
		err = errors.New(bs.String())
	}

	return err

}

//ApplicationRepo represents a repository abstraction for dealing with persistent Application instances.
type ApplicationRepo interface {
	CreateApplication(app *Application) error
	UpdateApplication(app *Application, subjectID string) error
	RetrieveApplication(clientID string, subjectID string, adminScope bool) (*Application, error)
	SystemRetrieveApplication(clientID string) (*Application, error)
	ListApplications(subjectID string, adminScope bool) ([]Application, error)
}

//NonOwnerUpdateError is used to discriminate general repo errors from security model violations
type NonOwnerUpdateError struct{}

//Error implements the Error interface for NonOwnerUpdateError
func (e NonOwnerUpdateError) Error() string {
	return "Non-owner attempted update"
}

//NoSuchApplicationError is used to discriminate updates of non-existent application from repo errors
type NoSuchApplicationError struct{}

//Error implements the Error interface for NoSuchApplicationError
func (e NoSuchApplicationError) Error() string {
	return "No such application to update"
}

//NotAuthorizedToReadApp is used to discriminate repo access errors from security model errors
type NotAuthorizedToReadApp struct{}

//Error implements the Error interface for NotAuthorizedToReadApp
func (e NotAuthorizedToReadApp) Error() string {
	return "Not authorized to read application definition"
}
