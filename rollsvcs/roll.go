package rollsvcs

import (
	"fmt"
	rollhttp "github.com/xtraclabs/roll/http"
	"github.com/xtraclabs/roll/repos"
	"github.com/xtraclabs/roll/repos/mdb"
	"github.com/xtraclabs/roll/roll"
	"log"
	"net/http"
)

func DefaultConfig() *roll.CoreConfig {
	return &roll.CoreConfig{
		DeveloperRepo:   repos.NewDynamoDevRepo(),
		ApplicationRepo: repos.NewDynamoAppRepo(),
		AdminRepo:       repos.NewDynamoAdminRepo(),
		SecretsRepo:     repos.NewVaultSecretsRepo(),
		IdGenerator:     new(roll.UUIDIdGenerator),
		Secure:          true,
	}
}

func DefaultUnsecureConfig() *roll.CoreConfig {
	return &roll.CoreConfig{
		DeveloperRepo:   repos.NewDynamoDevRepo(),
		ApplicationRepo: repos.NewDynamoAppRepo(),
		AdminRepo:       repos.NewDynamoAdminRepo(),
		SecretsRepo:     repos.NewVaultSecretsRepo(),
		IdGenerator:     new(roll.UUIDIdGenerator),
		Secure:          false,
	}
}

func MariaDBConfig() *roll.CoreConfig {
	return &roll.CoreConfig{
		AdminRepo:     mdb.NewMBDAdminRepo(),
		DeveloperRepo: mdb.NewMBDDevRepo(),
		IdGenerator:   new(roll.UUIDIdGenerator),
		Secure:        true,
	}
}

func RunRoll(port int, config *roll.CoreConfig) {
	core := roll.NewCore(config)
	log.Println("Starting roll - listening on port ", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), rollhttp.Handler(core))
}
