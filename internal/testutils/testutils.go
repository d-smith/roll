package testutils

import (
	"fmt"
	"github.com/xtraclabs/roll/roll"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"time"
)

func NewDevTestEmail() string {
	return fmt.Sprintf("%d@%d.net", time.Now().Unix(), time.Now().Unix())
}

func CreateNewTestDev() roll.Developer {
	return roll.Developer{
		FirstName: "test",
		LastName:  "test",
		Email:     NewDevTestEmail(),
	}
}

func URLGuard(url string) {
	const limit = 60
	var count int
	for {
		log.Info("check test endpoint availability")

		client := http.Client{}

		req,err := http.NewRequest("GET",url, nil)
		if err != nil {
			log.Info("Error creating request... bailing...", err.Error())
			return
		}

		req.Header.Set("X-Roll-Subject","rolltest")

		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK || count == limit {
			break
		}

		count += 1
		time.Sleep(1 * time.Second)
	}

	log.Info("acceptance tests ready for action boss")
}
