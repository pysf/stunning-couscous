package testutils

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/pysf/stunning-couscous/internal/db"
)

func CreateTestDatabase(t *testing.T) (testDB *sql.DB, tearDown func()) {
	testDBName := fmt.Sprintf("test_db_%v", rand.Intn(10000))

	psqDB, err := db.NewPostgreConnection("postgres")
	if err != nil {
		t.Fatalf("open db coection err: %s", err)
	}

	CreateDatabase(t, psqDB, testDBName)

	testDB, err = db.NewPostgreConnection(testDBName)
	if err != nil {
		t.Fatalf("open db coection err: %s", err)
	}

	tearDown = func() {
		if err := testDB.Close(); err != nil {
			t.Logf("failed to close connection to %v database: %s", testDBName, err)
		}

		_, err := psqDB.Exec(fmt.Sprintf("drop database %v", testDBName))
		if err != nil {
			t.Logf("teardown %v database failed: %s", testDBName, err)
		}
		psqDB.Close()
	}

	return testDB, tearDown
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
