package postgres

import (
	"database/sql"
	"io/ioutil"

	_ "github.com/lib/pq" // postgres driver
)

var tables []string = []string{
	"location",
	"account_user",
	"author",
	"work",
	"publication",
}

// DB wraps our SQL database.
type DB struct{ *sql.DB }

// Connect creates a pool of connections to the database
// and initializes the db on the receiver.
func (db *DB) Connect(params string) error {
	database, err := sql.Open("postgres", params)
	if err != nil {
		return err
	}

	if err = database.Ping(); err != nil {
		return err
	}

	db.DB = database
	return nil
}

// Disconnect closes the pool of connections to the database on the receiver.
func (db *DB) Disconnect() error {
	return db.Close()
}

// Init creates tables and populates them with data
// by running the scripts found at the given directory path.
func (db *DB) Init(dirPath string) error {
	for _, filePath := range getFilePaths(dirPath) {
		bytes, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		statements := string(bytes)
		if _, err := db.Exec(statements); err != nil {
			return err
		}
	}
	return nil
}

// Setup creates a new DB instance and returns it
// after successfully connecting and initializing the database.
func Setup(params string, sqlDirPath string) *DB {
	db := new(DB)

	if err := db.Connect(params); err != nil {
		panic(err)
	}

	if err := db.Init(sqlDirPath); err != nil {
		panic(err)
	}

	return db
}

func getFilePaths(dirPath string) []string {
	filePaths := []string{dirPath + "init.sql"}
	for _, table := range tables {
		filePath := dirPath + table + ".sql"
		filePaths = append(filePaths, filePath)
	}
	return filePaths
}
