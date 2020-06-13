package postgres

import (
	"database/sql"
	"io/ioutil"
	"path"
	"runtime"

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

// Query is used together with the Service.Query method to retrieve pre-defined queries.
type Query int

// Connect creates a pool of connections to the database and initializes the db on the receiver.
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

// Init creates tables and populates them with data.
func (db *DB) Init() error {
	for _, filePath := range getFilePaths() {
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

// Setup creates a new DB instance and returns it after successfully connecting
// and initializing the database. It also returns a closure function for easy cleanup.
func Setup(params string) (*DB, func()) {
	db := new(DB)

	if err := db.Connect(params); err != nil {
		panic(err)
	}

	if err := db.Init(); err != nil {
		panic(err)
	}

	return db, func() { db.Disconnect() }
}

func getFilePaths() []string {
	_, file, _, _ := runtime.Caller(0)
	dirPath := path.Join(path.Dir(file), "../sql")

	filePaths := []string{path.Join(dirPath, "init.sql")}
	for _, table := range tables {
		filePath := path.Join(dirPath, table+".sql")
		filePaths = append(filePaths, filePath)
	}
	return filePaths
}
