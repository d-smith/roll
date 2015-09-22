package roll

import (
	"bytes"
	"errors"
	"regexp"
)

//Developer represents the data associated with a Developer
type Developer struct {
	FirstName string
	LastName  string
	Email     string
	ID        string
}

var validEmail = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
var firstName = regexp.MustCompile(`^[a-zA-Z]+$`)
var lastName = regexp.MustCompile(`^([a-zA-Z'-]\s*)+$`)

//ValidateEmail returns true when given a valid email address
func ValidateEmail(email string) bool {
	return validEmail.MatchString(email)
}

func (d *Developer) validateEmail() bool {
	return validEmail.MatchString(d.Email)
}

func (d *Developer) validateFirstName() bool {
	return firstName.MatchString(d.FirstName)
}

func (d *Developer) validateLastName() bool {
	return lastName.MatchString(d.LastName)
}

//Validate tests all the field in Developer to make sure they contain valid content. An error
//is returned if they don't.
func (d *Developer) Validate() error {
	var valid = true
	var err error

	bs := bytes.NewBufferString("Fields with invalid content: ")

	if !d.validateEmail() {
		valid = false
		bs.WriteString("Email ")
	}

	if !d.validateFirstName() {
		valid = false
		bs.WriteString("FirstName ")
	}

	if !d.validateLastName() {
		valid = false
		bs.WriteString("LastName ")
	}

	if !valid {
		err = errors.New(bs.String())
	}

	return err
}

//DeveloperRepo represents a repository abstraction for dealing with persistent Developer instances.
type DeveloperRepo interface {
	RetrieveDeveloper(email string) (*Developer, error)
	StoreDeveloper(*Developer) error
}
