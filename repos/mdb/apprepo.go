package mdb

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/xtraclabs/roll/dbutil"
	"github.com/xtraclabs/roll/repos"
	"github.com/xtraclabs/roll/roll"
	"github.com/xtraclabs/roll/secrets"
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
	//Generate a client secret as needed
	if app.ClientSecret == "" {
		clientSecret, err := secrets.GenerateClientSecret()
		if err != nil {
			return err
		}
		app.ClientSecret = clientSecret
	}

	//Check JWT flow parts are ok
	if err := repos.CheckJWTCertParts(app); err != nil {
		return err
	}

	//Insert the app
	const appSql = `insert into rolldb.application(applicationName, clientId, clientSecret, developerEmail, developerId, loginProvider,
	redirectUri,jwtFlowAudience, jwtFlowIssuer, jwtFlowPublicKey) values(?,?,?,?,?,?,?,?,?,?)
	`
	stmt, err := ar.db.Prepare(appSql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		app.ApplicationName,
		app.ClientID,
		app.ClientSecret,
		app.DeveloperEmail,
		app.DeveloperID,
		app.LoginProvider,
		app.RedirectURI,
		app.JWTFlowAudience,
		app.JWTFlowIssuer,
		app.JWTFlowPublicKey,
	)

	if err != nil {
		log.Println(err)
		sqlErr := err.(*mysql.MySQLError)
		switch sqlErr.Number {
		case 1062:
			log.Println("Duplicate app definition found")
			return repos.NewDuplicationAppdefError(app.ApplicationName, app.DeveloperEmail)
		default:
			return err
		}
	}

	return nil
}

func applyUpdatesWithJWTColumns(db *sql.DB, app *roll.Application) error {
	//Check JWT flow parts are ok
	if err := repos.CheckJWTCertParts(app); err != nil {
		return err
	}

	const updateSql = `
	update application set loginProvider=?, redirectUri=?,jwtFlowPublicKey=?,jwtFlowIssuer=?,
	jwtFlowAudience=?,applicationName=? where clientId=?
	`
	stmt, err := db.Prepare(updateSql)
	if err != nil {
		return nil
	}

	defer stmt.Close()

	_, err = stmt.Exec(app.LoginProvider, app.RedirectURI, app.JWTFlowPublicKey, app.JWTFlowIssuer,
		app.JWTFlowAudience, app.ApplicationName, app.ClientID)
	return err

}

func applyUpdate(db *sql.DB, app *roll.Application) error {
	const updateSql = `
	update application set loginProvider=?, redirectUri=?,applicationName=? where clientId=?
	`
	stmt, err := db.Prepare(updateSql)
	if err != nil {
		return nil
	}

	defer stmt.Close()

	_, err = stmt.Exec(app.LoginProvider, app.RedirectURI, app.ApplicationName, app.ClientID)
	return err
}

//to the existing things, specifically for the jwt parts?
func (ar *MariaDBAppRepo) UpdateApplication(app *roll.Application, subjectID string) error {
	storedApp, err := ar.SystemRetrieveApplication(app.ClientID)
	if err != nil {
		log.Println("Error retrieving app to verify ownership")
		return err
	}

	if storedApp == nil {
		log.Println("Application to update does not exist")
		return roll.NoSuchApplicationError{}
	}

	if storedApp.DeveloperID != subjectID {
		log.Println("Application updater does not own app (" + subjectID + ")")
		return roll.NonOwnerUpdateError{}
	}

	log.Println("Updating", app.ClientID, "owned by", app.DeveloperID)

	if app.JWTFlowPublicKey != "" {
		err = applyUpdatesWithJWTColumns(ar.db, app)
	} else {
		err = applyUpdate(ar.db, app)
	}

	return err

}

func (ar *MariaDBAppRepo) delete(app *roll.Application) error {
	db := ar.db

	stmt, err := db.Prepare("delete from application where applicationName = ? and developerEmail = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(app.ApplicationName, app.DeveloperEmail)
	if err != nil {
		return err
	}

	return nil
}

func (ar *MariaDBAppRepo) RetrieveAppByNameAndDevEmail(appName, email string) (*roll.Application, error) {
	const appSql = `
	select applicationName, clientId, clientSecret, developerEmail, developerId, loginProvider,
	redirectUri,jwtFlowAudience, jwtFlowIssuer, jwtFlowPublicKey from application where applicationName = ?
	and developerEmail = ?
	`

	var app roll.Application
	err := ar.db.QueryRow(appSql,
		appName, email).Scan(
		&app.ApplicationName, &app.ClientID, &app.ClientSecret, &app.DeveloperEmail, &app.DeveloperID, &app.LoginProvider,
		&app.RedirectURI, &app.JWTFlowAudience, &app.JWTFlowIssuer, &app.JWTFlowPublicKey,
	)

	if err != nil {
		return nil, err
	}

	return &app, nil
}

func (ar *MariaDBAppRepo) RetrieveApplication(clientID string, subjectID string, adminScope bool) (*roll.Application, error) {
	app, err := ar.SystemRetrieveApplication(clientID)
	if err != nil {
		return nil, err
	}

	if !adminScope && app.DeveloperID != subjectID {
		return nil, roll.NotAuthorizedToReadApp{}
	}

	return app, nil
}

//SystemRetrieveApplication is used for system level access of application records where the user
//security model does not need to be applied.
func (ar *MariaDBAppRepo) SystemRetrieveApplication(clientID string) (*roll.Application, error) {
	const appSql = `
	select applicationName, clientId, clientSecret, developerEmail, developerId, loginProvider,
	redirectUri,jwtFlowAudience, jwtFlowIssuer, jwtFlowPublicKey from application where clientId = ?
	`

	log.Println("Looking up app for", clientID)
	var app roll.Application
	err := ar.db.QueryRow(appSql,
		clientID).Scan(
		&app.ApplicationName, &app.ClientID, &app.ClientSecret, &app.DeveloperEmail, &app.DeveloperID, &app.LoginProvider,
		&app.RedirectURI, &app.JWTFlowAudience, &app.JWTFlowIssuer, &app.JWTFlowPublicKey,
	)

	return &app, err
}

func (ar *MariaDBAppRepo) SystemRetrieveApplicationByJWTFlowAudience(audience string) (*roll.Application, error) {
	const appSql = `
	select applicationName, clientId, clientSecret, developerEmail, developerId, loginProvider,
	redirectUri,jwtFlowAudience, jwtFlowIssuer, jwtFlowPublicKey from application where jwtFlowAudience = ?
	`

	var app roll.Application
	err := ar.db.QueryRow(appSql,
		audience).Scan(
		&app.ApplicationName, &app.ClientID, &app.ClientSecret, &app.DeveloperEmail, &app.DeveloperID, &app.LoginProvider,
		&app.RedirectURI, &app.JWTFlowAudience, &app.JWTFlowIssuer, &app.JWTFlowPublicKey,
	)

	return &app, err
}

func (ar *MariaDBAppRepo) ListApplications(subjectID string, adminScope bool) ([]roll.Application, error) {
	var rows *sql.Rows
	var err error

	if adminScope == true {
		const adminScopeSelect = `
		select applicationName, clientId, developerEmail, developerId, loginProvider,
		redirectUri,jwtFlowAudience, jwtFlowIssuer, jwtFlowPublicKey from application
		`

		rows, err = ar.db.Query(adminScopeSelect)
	} else {
		const nonAdminSelect = `
		select applicationName, clientId, developerEmail, developerId, loginProvider,
		redirectUri,jwtFlowAudience, jwtFlowIssuer, jwtFlowPublicKey from application where developerId = ?
		`

		rows, err = ar.db.Query(nonAdminSelect, subjectID)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var apps []roll.Application
	for rows.Next() {
		app := new(roll.Application)
		err = rows.Scan(&app.ApplicationName,
			&app.ClientID,
			&app.DeveloperEmail,
			&app.DeveloperID,
			&app.LoginProvider,
			&app.RedirectURI,
			&app.JWTFlowAudience,
			&app.JWTFlowIssuer,
			&app.JWTFlowPublicKey,
		)

		if err != nil {
			return nil, err
		}

		apps = append(apps, *app)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return apps, nil
}
