package partner_test

import (
	"context"
	"testing"

	"github.com/pysf/stunning-couscous/internal/partner"
)

func TestPartnerPepoFindBestMatch_ValidateDistance(t *testing.T) {
	repo, err := partner.NewPartnerRepo()
	if err != nil {
		t.Fatalf("NewPartnerRepo() err = %v", err)
	}

	arg1 := partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}
	arg2 := "wood"

	got, err := repo.FindBestMatch(context.Background(), arg1, arg2)
	if err != nil {
		t.Fatalf("NewPartnerRepo() err = %v", err)
	}

	for _, p := range got {
		if p.Distance > p.OperatingRadius {
			t.Errorf("FindBestMatch() = %v distance is gt operatingRange %v ", p.Distance, p.OperatingRadius)
		}
	}

}
