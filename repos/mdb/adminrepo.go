package mdb

import (
	"database/sql"
	"github.com/xtraclabs/roll/dbutil"
	"log"
)

type MBDAdminRepo struct {
	db *sql.DB
}

func NewMBDAdminRepo() *MBDAdminRepo {
	//If we error out, there nothing we can do to recover, so we're done.
	db, err := dbutil.CreateMariaDBSqlDB()
	if err != nil {
		log.Fatal("Error prepping for MariaDB connection", err.Error())
	}
	return &MBDAdminRepo{
		db: db,
	}
}

func (ar *MBDAdminRepo) IsAdmin(subject string) (bool, error) {
	var count int

	err := ar.db.QueryRow("select count(*) from admin where name = ?", subject).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
