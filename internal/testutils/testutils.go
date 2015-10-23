package testutils

import (
	"fmt"
	"time"
	"github.com/xtraclabs/roll/roll"
	"log"
	"net/http"
)

func NewDevTestEmail() string {
	return fmt.Sprintf("%d@%d.net", time.Now().Unix(), time.Now().Unix())
}

func CreateNewTestDev() roll.Developer {
	return roll.Developer{
		FirstName: "test",
		LastName: "test",
		Email: NewDevTestEmail(),
	}
}

func URLGuard(url string) {
	const limit = 60
	var count int
	for {
		log.Println("check test endpoint availability")
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK || count == limit {
			break
		}

		count += 1
		time.Sleep(1*time.Second)
	}

	log.Println("acceptance tests ready for action boss")
}
