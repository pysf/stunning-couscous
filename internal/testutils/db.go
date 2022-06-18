package testutils

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
)

func CreateTestDatabase(t *testing.T) (testDB *sql.DB, tearDown func()) {
	testDBName := fmt.Sprintf("test_db_%v", rand.Intn(10000))

	db := DBConnection(t, "postgres")
	CreateDatabase(t, db, testDBName)

	testDB = DBConnection(t, testDBName)

	tearDown = func() {
		if err := testDB.Close(); err != nil {
			t.Logf("failed to close connection to test database: %s", err)
		}

		_, err := db.Exec(fmt.Sprintf("drop database %v", testDBName))
		if err != nil {
			t.Logf("teardown test database failed: %s", err)
		}
		db.Close()
	}

	return testDB, tearDown
}

func DBConnection(t *testing.T, dbName string) *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRESQL_HOST"),
		os.Getenv("POSTGRESQL_PORT"),
		os.Getenv("POSTGRESQL_USERNAME"),
		os.Getenv("POSTGRESQL_PASSWORD"),
		dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		t.Fatalf("failed to open db connection %s", err)
	}

	if err = db.Ping(); err != nil {
		t.Fatalf("failed to ping connection %s", err)
	}

	return db
}

func CreateDatabase(t *testing.T, db *sql.DB, dbname string) {
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %v", dbname))
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return
		}
		t.Errorf("failed to create test database %s", err)
	}
}
