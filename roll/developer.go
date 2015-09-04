package roll

import "regexp"

//Developer represents the data associated with a Developer
type Developer struct {
	FirstName string
	LastName  string
	Email     string
	ID        string
}

//ValidateEmail validates a string that is to be treated as an email address
func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

//DeveloperRepo represents a repository abstraction for dealing with persistent Developer instances.
type DeveloperRepo interface {
	RetrieveDeveloper(email string) (*Developer, error)
	StoreDeveloper(*Developer) error
}
