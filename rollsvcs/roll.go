package rollsvcs

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	rollhttp "github.com/xtraclabs/roll/http"
	"github.com/xtraclabs/roll/repos"
	"github.com/xtraclabs/roll/repos/mdb"
	"github.com/xtraclabs/roll/roll"
	secretsrepos "github.com/xtraclabs/rollsecrets/repos"
	rolltoken "github.com/xtraclabs/rollsecrets/token"
	"net/http"
)

func DefaultConfig() *roll.CoreConfig {
	return &roll.CoreConfig{
		DeveloperRepo:   repos.NewDynamoDevRepo(),
		ApplicationRepo: repos.NewDynamoAppRepo(),
		AdminRepo:       repos.NewDynamoAdminRepo(),
		SecretsRepo:     secretsrepos.NewVaultSecretsRepo(),
		IdGenerator:     new(rolltoken.UUIDIdGenerator),
		Secure:          true,
	}
}

func DefaultUnsecureConfig() *roll.CoreConfig {
	return &roll.CoreConfig{
		DeveloperRepo:   repos.NewDynamoDevRepo(),
		ApplicationRepo: repos.NewDynamoAppRepo(),
		AdminRepo:       repos.NewDynamoAdminRepo(),
		SecretsRepo:     secretsrepos.NewVaultSecretsRepo(),
		IdGenerator:     new(rolltoken.UUIDIdGenerator),
		Secure:          false,
	}
}

func MariaDBUnsecureConfig() *roll.CoreConfig {
	return &roll.CoreConfig{
		AdminRepo:       mdb.NewMBDAdminRepo(),
		DeveloperRepo:   mdb.NewMBDDevRepo(),
		ApplicationRepo: mdb.NewMBDAppRepo(),
		SecretsRepo:     secretsrepos.NewVaultSecretsRepo(),
		IdGenerator:     new(rolltoken.UUIDIdGenerator),
		Secure:          false,
	}
}

func MariaDBSecureConfig() *roll.CoreConfig {
	return &roll.CoreConfig{
		AdminRepo:       mdb.NewMBDAdminRepo(),
		DeveloperRepo:   mdb.NewMBDDevRepo(),
		ApplicationRepo: mdb.NewMBDAppRepo(),
		SecretsRepo:     secretsrepos.NewVaultSecretsRepo(),
		IdGenerator:     new(rolltoken.UUIDIdGenerator),
		Secure:          true,
	}
}

func RunRoll(port int, config *roll.CoreConfig) {
	core := roll.NewCore(config)
	log.Info("Starting roll - listening on port ", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), rollhttp.Handler(core))
}
