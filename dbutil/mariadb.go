package dbutil

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func UseMariaDB() bool {
	return os.Getenv("ROLL_DBADDRESS") != ""
}

func CreateMariaDBSqlDB() (*sql.DB, error) {
	rollUser := os.Getenv("ROLL_DBUSER")
	rollPassword := os.Getenv("ROLL_DBPASSWORD")
	rollAddress := os.Getenv("ROLL_DBADDRESS")

	connectString := fmt.Sprintf("%s:%s@tcp(%s)/rolldb", rollUser, rollPassword, rollAddress)

	return sql.Open("mysql", connectString)
}
