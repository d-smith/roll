package developers


import (
	. "github.com/lsegal/gucumber"
	"github.com/xtraclabs/roll/roll"
	"time"
	"fmt"
	rollhttp "github.com/xtraclabs/roll/http"
	"github.com/stretchr/testify/assert"
	"net/http"
	"log"
)

func devEmail() string {
	return fmt.Sprintf("%d@%d.net", time.Now().Unix(), time.Now().Unix())
}

func init() {

	var newDev roll.Developer
	var malformed roll.Developer

	Given(`^A developer who registers on the portal$`, func() {
		log.Println("Create dev with malformed email")
		newDev = roll.Developer{
			FirstName: "test",
			LastName: "test",
			Email: devEmail(),
		}
	})

	And(`^They have not registered before$`, func() {
		//Assumed
	})

	And(`^Their data is formatted correctly$`, func() {
		//Ensured in the above given
	})

	Then(`^They are added to the portal successfully$`, func() {
		resp := rollhttp.TestHTTPPut(T, "http://localhost:3000/v1/developers/" + newDev.Email, newDev)
		assert.Equal(T, http.StatusNoContent, resp.StatusCode)
	})

	Given(`^a developer who registers on the portal$`, func() {
		log.Println("Create dev with malformed email")
		malformed = roll.Developer{
			FirstName: "test",
			LastName: "test",
			Email: "no-good",
		}
	})

	And(`^They provide a malformed email$`, func() {
		//Ensured above
	})

	Then(`^An error is returned with StatusBadRequest$`, func() {
		log.Println("Add dev with malformed email", malformed)
		resp := rollhttp.TestHTTPPut(T, "http://localhost:3000/v1/developers/" + malformed.Email, malformed)
		assert.Equal(T, http.StatusBadRequest, resp.StatusCode)
	})
	
}