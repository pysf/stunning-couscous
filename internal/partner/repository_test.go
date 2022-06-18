package partner_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/pysf/stunning-couscous/internal/bulkgen"
	"github.com/pysf/stunning-couscous/internal/partner"
	"github.com/pysf/stunning-couscous/internal/testutils"
)

func TestPartnerRepoFindBestMatch_ValidateDistance(t *testing.T) {
	db, tearDown := setupDB(t)
	defer tearDown()

	baseLocation := partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}
	seedTestPartners(t, db, baseLocation, 300)

	repo := partner.PartnerRepo{
		DB: db,
	}

	arg1 := partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}
	arg2 := "wood"

	got, err := repo.FindBestMatch(context.Background(), arg1, arg2)
	if err != nil {
		t.Fatalf("NewPartnerRepo() err = %s", err)
	}

	for _, p := range got {
		if p.Distance > p.OperatingRadius {
			t.Fatalf("FindBestMatch() = %v ; want distance lt %v ", p.Distance, p.OperatingRadius)
		}
	}

}

func TestPartnerRepoFindBestMatch_ValidateExperience(t *testing.T) {
	db, tearDown := setupDB(t)
	defer tearDown()

	baseLocation := partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}
	seedTestPartners(t, db, baseLocation, 300)

	repo := partner.PartnerRepo{
		DB: db,
	}

	arg1 := partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}
	arg2 := "tiles"

	got, err := repo.FindBestMatch(context.Background(), arg1, arg2)
	if err != nil {
		t.Fatalf("NewPartnerRepo() err = %v", err)
	}

	for _, p := range got {
		if !Contains(p.Experiences, arg2) {
			t.Fatalf("FindBestMatch() = %v , want %q experience ", p.Experiences, arg2)
		}
	}

}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func setupDB(t *testing.T) (*sql.DB, func()) {
	db, tearDown := testutils.CreateTestDatabase(t)
	partner.ApplySchema(db)
	return db, tearDown
}

func seedTestPartners(t *testing.T, db *sql.DB, loc partner.Location, size int) {

	locations := bulkgen.GenerateRandomLocations(loc, size)
	partners := bulkgen.GeneratePartner(locations)

	partnerRepo := partner.PartnerRepo{
		DB: db,
	}

	partnerRepo.BulkImport(partners)

}
