package mdb

import (
	"database/sql"
	"github.com/xtraclabs/roll/dbutil"
	"github.com/xtraclabs/roll/roll"
	"log"
)

type MariaDBAppRepo struct {
	db *sql.DB
}

func NewMBDAppRepo() *MariaDBAppRepo {
	//If we error out, there nothing we can do to recover, so we're done.
	db, err := dbutil.CreateMariaDBSqlDB()
	if err != nil {
		log.Fatal("Error prepping for MariaDB connection", err.Error())
	}
	return &MariaDBAppRepo{
		db: db,
	}
}

func (ar *MariaDBAppRepo) CreateApplication(app *roll.Application) error {
	return nil
}

func (ar *MariaDBAppRepo) UpdateApplication(app *roll.Application, subjectID string) error {
	return nil
}

func (ar *MariaDBAppRepo) RetrieveAppByNameAndDevEmail(appName, email string) (*roll.Application, error) {
	return nil, nil
}

func (ar *MariaDBAppRepo) RetrieveApplication(clientID string, subjectID string, adminScope bool) (*roll.Application, error) {
	return nil, nil
}

func (ar *MariaDBAppRepo) SystemRetrieveApplication(clientID string) (*roll.Application, error) {
	return nil, nil
}

func (ar *MariaDBAppRepo) SystemRetrieveApplicationByJWTFlowAudience(audience string) (*roll.Application, error) {
	return nil, nil
}

func (ar *MariaDBAppRepo) ListApplications(subjectID string, adminScope bool) ([]roll.Application, error) {
	return nil, nil
}
