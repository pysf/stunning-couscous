package bulkgen_test

import (
	"reflect"
	"testing"

	"github.com/pysf/stunning-couscous/internal/bulkgen"
	"github.com/pysf/stunning-couscous/internal/partner"
)

func TestGenerateRandomLocations_Diversity(t *testing.T) {
	baseLocation := partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}

	size := 100
	got := bulkgen.GenerateRandomLocations(baseLocation, size)

	var current partner.Location
	for _, l := range got {
		if reflect.DeepEqual(current, l) {
			t.Fatalf("GenerateRandomLocations() = %v, is duplicate, must be unique %v = %v", l, l, current)
		}
		current = l
	}
}

func TestGenerateRandomLocations_Size(t *testing.T) {

	baseLocation := partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}

	size := 100
	got := bulkgen.GenerateRandomLocations(baseLocation, size)

	if len(got) < size {
		t.Fatalf("GenerateRandomLocations() size= %v, must be gt %v", len(got), size)
	}

}

func TestGeneratePartner(t *testing.T) {

	wantedSize := 18
	baseLocation := partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}
	locations := bulkgen.GenerateRandomLocations(baseLocation, wantedSize)

	partners := bulkgen.GeneratePartner(locations)

	if len(partners) != wantedSize {
		t.Errorf("GeneratePartner() = %v, wanted size %v", len(partners), wantedSize)
	}

	for _, got := range partners {

		if got.Longitude == baseLocation.Longitude && got.Latitude == baseLocation.Latitude {
			t.Errorf("GeneratePartner() = (%v,%v), same as base  (%v,%v)", got.Latitude, got.Longitude, baseLocation.Latitude, baseLocation.Longitude)
		}

		if len(got.Experiences) == 0 {
			t.Errorf("GeneratePartner() , Experiences is not filled %v", len(got.Experiences))
		}

		if got.Rating == 0 {
			t.Errorf("GeneratePartner() , Rating can not be %v", got.Rating)
		}

		if got.OperatingRadius == 0 {
			t.Errorf("GeneratePartner() , OperatingRadius can not be %v", got.OperatingRadius)
		}
	}

}
