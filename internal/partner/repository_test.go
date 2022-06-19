package partner_test

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/lib/pq"
	"github.com/pysf/stunning-couscous/internal/bulkgen"
	"github.com/pysf/stunning-couscous/internal/partner"
	"github.com/pysf/stunning-couscous/internal/testutils"
)

func TestGetPartner(t *testing.T) {
	db, tearDown := testutils.SetupDB(t)
	defer tearDown()

	repo := partner.PartnerRepo{
		DB: db,
	}

	want := &partner.Partner{
		ID:              101,
		Rating:          4,
		Distance:        10,
		OperatingRadius: 8,
		Experiences:     []string{"wood", "carpet"},
		Location: partner.Location{
			Latitude:  52.528971849007036,
			Longitude: 13.430548464498173,
		},
	}
	if _, err := repo.DB.Exec(`INSERT INTO  partner ("id", "experiences", "operatingradius", "rating", "location") VALUES($1, $2, $3, $4, POINT( $5, $6 ) )`,
		want.ID, pq.Array(want.Experiences), want.OperatingRadius, want.Rating, want.Latitude, want.Longitude); err != nil {
		t.Errorf("failed to add partner %s", err)
	}

	got, err := repo.GetPartner(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("GetPartner() err = %s", err)
	}

	if reflect.DeepEqual(want, got) {
		t.Fatalf("GetPartner() = %v; want %v", got, want)
	}

}

func TestPartnerRepoFindBestMatch_ValidateDistance(t *testing.T) {
	db, tearDown := testutils.SetupDB(t)
	defer tearDown()

	baseLocation := partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}
	seedTestPartners(t, db, baseLocation, 300)

	repo := partner.PartnerRepo{
		DB: db,
	}

	arg1 := baseLocation
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
	db, tearDown := testutils.SetupDB(t)
	defer tearDown()

	baseLocation := partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}
	seedTestPartners(t, db, baseLocation, 300)

	repo := partner.PartnerRepo{
		DB: db,
	}

	arg1 := baseLocation
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

func seedTestPartners(t *testing.T, db *sql.DB, loc partner.Location, size int) {

	locations := bulkgen.GenerateRandomLocations(loc, size)
	partners := bulkgen.GeneratePartner(locations)

	partnerRepo := partner.PartnerRepo{
		DB: db,
	}

	partnerRepo.BulkImport(partners)

}
