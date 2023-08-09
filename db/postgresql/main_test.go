package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pawpaw2022/simplebank/util"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	var err error

	// Load config
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config: %w", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: %w", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
