package db

import (
	"database/sql"
	"io/ioutil"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres driver
)

var (
	configPath string   = "config/config.env"
	sqlFiles   []string = []string{"db/sql/schema.sql", "db/sql/init.sql"}
)

// Connect creates and returns a pool of connections to the database.
func Connect() (*sql.DB, error) {
	err := godotenv.Load(configPath)
	if err != nil {
		return nil, err
	}

	connectionString := os.Getenv("CONNECTION_STRING")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Disconnect closes the pool of connections to the given database.
func Disconnect(db *sql.DB) {
	db.Close()
}

// Init creates tables by running the appropriate SQL scripts
// and also inserts existing data that we have into the tables.
func Init(database *sql.DB) error {
	for _, path := range sqlFiles {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		statements := string(bytes)
		if _, err := database.Exec(statements); err != nil {
			return err
		}
	}
	return nil
}
