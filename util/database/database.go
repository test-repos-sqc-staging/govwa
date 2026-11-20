package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/moshelior/govwa/util/config"
	"log"
	"regexp"
)

// validDbName only allows alphanumeric characters and underscores, preventing
// SQL injection through the database name since identifiers cannot be
// parameterized by the driver.
var validDbName = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func Connect() (*sql.DB, error) {

	config := config.LoadConfig()

	var dsn string
	var db *sql.DB

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/", config.User, config.Password, config.Sqlhost, config.Sqlport)
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}
	if !validDbName.MatchString(config.Dbname) {
		return nil, fmt.Errorf("invalid database name: %s", config.Dbname)
	}
	// config.Dbname is validated above against validDbName, so it cannot
	// be used to alter the structure of this statement.
	createDbQuery := "CREATE DATABASE IF NOT EXISTS " + config.Dbname
	_, err = db.Exec(createDbQuery) // NOSONAR

	if err != nil {
		return nil, err
	} else {

		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Password, config.Sqlhost, config.Sqlport, config.Dbname)
		db, err = sql.Open("mysql", dsn)

		if err != nil {
			return nil, err

		}
	}

	return db, nil
}

var DB *sql.DB

func CheckDatabase() (bool, error) {

	/* this function use to check if no database selected and will redirect to setup page */

	DB, err := Connect()
	if err != nil {
		log.Printf("Connection Error %s ", err.Error())
	}

	const (
		checksql = `SELECT 1  FROM Users limit 1` //this will check if Table dbname.Users exist otherways will redirect to setup page
	)
	result, err := DB.Exec(checksql)

	if err != nil {
		log.Println(err.Error())
		return false, err
	}
	if result == nil {
		return false, err
	}
	log.Println(result)
	return true, nil
}
