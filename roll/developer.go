package roll

import "regexp"

type Developer struct {
	FirstName string
	LastName  string
	Email     string
	Id        string
}

func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

type DeveloperRepo interface {
	RetrieveDeveloper(email string) (*Developer, error)
	StoreDeveloper(*Developer) error
}
