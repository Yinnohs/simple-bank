package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/yinnohs/simple-bank/util"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://yinnohs:1234@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal(" cannot load configuration", err)
		config = util.AppConfig{
			DbDriver: dbDriver,
			DbSource: dbSource,
		}
	}

	testDB, err = sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatal("Caanot connect to db: ", err.Error())
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
