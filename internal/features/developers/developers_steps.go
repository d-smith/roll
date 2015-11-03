package developers

import (
	. "github.com/lsegal/gucumber"
	"github.com/stretchr/testify/assert"
	rollhttp "github.com/xtraclabs/roll/http"
	"github.com/xtraclabs/roll/internal/testutils"
	"github.com/xtraclabs/roll/roll"
	"log"
	"net/http"
)

func init() {

	var newDev roll.Developer
	var malformed roll.Developer

	Before("@devtests", func() {
		testutils.URLGuard("http://localhost:3000/v1/developers")
	})

	Given(`^A developer who registers on the portal$`, func() {
		log.Println("Create dev with malformed email")
		newDev = testutils.CreateNewTestDev()
	})

	And(`^They have not registered before$`, func() {
		//Assumed
	})

	And(`^Their data is formatted correctly$`, func() {
		//Ensured in the above given
	})

	Then(`^They are added to the portal successfully$`, func() {
		resp := rollhttp.TestHTTPPutWithRollSubject(T, "http://localhost:3000/v1/developers/"+newDev.Email, newDev)
		assert.Equal(T, http.StatusNoContent, resp.StatusCode)
	})

	Given(`^a developer who registers on the portal$`, func() {
		log.Println("Create dev with malformed email")
		malformed = roll.Developer{
			FirstName: "test",
			LastName:  "test",
			Email:     "no-good",
		}
	})

	And(`^They provide a malformed email$`, func() {
		//Ensured above
	})

	Then(`^An error is returned with StatusBadRequest$`, func() {
		log.Println("Add dev with malformed email", malformed)
		resp := rollhttp.TestHTTPPutWithRollSubject(T, "http://localhost:3000/v1/developers/"+malformed.Email, malformed)
		assert.Equal(T, http.StatusBadRequest, resp.StatusCode)
	})

}
