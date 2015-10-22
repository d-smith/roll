package rollsvcs

import (
	"fmt"
	rollhttp "github.com/xtraclabs/roll/http"
	"github.com/xtraclabs/roll/repos"
	"github.com/xtraclabs/roll/roll"
	"log"
	"net/http"
)

func DefaultConfig() *roll.CoreConfig {
	return &roll.CoreConfig{
		DeveloperRepo:   repos.NewDynamoDevRepo(),
		ApplicationRepo: repos.NewDynamoAppRepo(),
		SecretsRepo:     repos.NewVaultSecretsRepo(),
		IdGenerator: new(roll.UUIDIdGenerator),
	}
}

func RunRoll(port int, config *roll.CoreConfig) {
	core := roll.NewCore(config)
	log.Println("Starting roll - listening on port ", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), rollhttp.Handler(core))
}
