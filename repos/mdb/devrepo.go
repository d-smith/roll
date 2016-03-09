package mdb

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	"github.com/xtraclabs/roll/dbutil"
	"github.com/xtraclabs/roll/roll"
)

type MBDDevRepo struct {
	db *sql.DB
}

func NewMBDDevRepo() *MBDDevRepo {
	//If we error out, there nothing we can do to recover, so we're done.
	db, err := dbutil.CreateMariaDBSqlDB()
	if err != nil {
		log.Fatal("Error prepping for MariaDB connection", err.Error())
	}
	return &MBDDevRepo{
		db: db,
	}
}

func doAdminDevQuery(email string, db *sql.DB) (*roll.Developer, error) {
	var dev roll.Developer
	err := db.QueryRow("select firstName, lastName, email, id from developer where email = ?", email).Scan(
		&dev.FirstName, &dev.LastName, &dev.Email, &dev.ID,
	)

	return &dev, err
}

func doUserDevQuery(email, subject string, db *sql.DB) (*roll.Developer, error) {
	var dev roll.Developer
	err := db.QueryRow("select firstName, lastName, email, id from developer where email = ? and id = ?",
		email, subject).Scan(
		&dev.FirstName, &dev.LastName, &dev.Email, &dev.ID,
	)

	return &dev, err
}

func (dr *MBDDevRepo) RetrieveDeveloper(email string, subjectID string, adminScope bool) (*roll.Developer, error) {
	switch adminScope {
	case true:
		return doAdminDevQuery(email, dr.db)
	default:
		return doUserDevQuery(email, subjectID, dr.db)
	}
}

func (dr *MBDDevRepo) StoreDeveloper(dev *roll.Developer) error {
	db := dr.db

	stmt, err := db.Prepare("insert into developer(id,email, firstName, lastName) values (?,?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(dev.ID, dev.Email, dev.FirstName, dev.LastName)
	if err != nil {
		return err
	}

	return nil
}

func processListRows(rows *sql.Rows) ([]roll.Developer, error) {
	defer rows.Close()

	var devs []roll.Developer
	for rows.Next() {
		dev := new(roll.Developer)
		err := rows.Scan(&dev.FirstName, &dev.LastName, &dev.ID, &dev.Email)
		if err != nil {
			return nil, err
		}
		devs = append(devs, *dev)
	}

	err := rows.Err()
	if err != nil {
		return nil, err
	}

	return devs, nil
}

func adminListDevs(db *sql.DB) ([]roll.Developer, error) {
	rows, err := db.Query("select firstName, lastName, id, email from developer")
	if err != nil {
		return nil, err
	}

	return processListRows(rows)

}

func (dr *MBDDevRepo) deleteDeveloper(email string) error {
	db := dr.db

	stmt, err := db.Prepare("delete from developer where email = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(email)
	if err != nil {
		return err
	}

	return nil
}

func devListDevs(db *sql.DB, subject string) ([]roll.Developer, error) {
	rows, err := db.Query("select firstName, lastName, id, email from developer where id = ?", subject)
	if err != nil {
		return nil, err
	}
	return processListRows(rows)
}

func (dr *MBDDevRepo) ListDevelopers(subjectID string, adminScope bool) ([]roll.Developer, error) {
	switch adminScope {
	case true:
		return adminListDevs(dr.db)
	default:
		return devListDevs(dr.db, subjectID)
	}

}
