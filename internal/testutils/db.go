package testutils

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/pysf/stunning-couscous/internal/bulkgen"
	"github.com/pysf/stunning-couscous/internal/db"
	"github.com/pysf/stunning-couscous/internal/partner"
)

func CreateTestDatabase(t *testing.T) (testDB *sql.DB, tearDown func()) {
	rand.Seed(time.Now().UnixNano())

	testDBName := fmt.Sprintf("test_db_%v", rand.Intn(100000))

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

func SetupDB(t *testing.T) (*sql.DB, func()) {
	db, tearDown := CreateTestDatabase(t)
	partner.ApplySchema(db)
	return db, tearDown
}

func SeedTestPartners(t *testing.T, db *sql.DB, loc partner.Location, size int) {

	locations := bulkgen.GenerateRandomLocations(loc, size)
	partners := bulkgen.GeneratePartner(locations)

	partnerRepo := partner.PartnerRepo{
		DB: db,
	}

	partnerRepo.BulkImport(partners)

}
