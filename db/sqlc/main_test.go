package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:postgres@localhost:54320/cheapbank?sslmode=disable"
)

func TestMain(m *testing.M) {  
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(testDB)
	
	code := m.Run()
	
	testDB.Close()  
	os.Exit(code)
}


