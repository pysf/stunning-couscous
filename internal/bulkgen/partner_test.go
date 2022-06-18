package bulkgen_test

import (
	"reflect"
	"testing"

	"github.com/pysf/stunning-couscous/internal/bulkgen"
	"github.com/pysf/stunning-couscous/internal/partner"
)

func TestGenerateRandomLocations_diversity(t *testing.T) {
	baseLocation := partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}

	size := 100000
	got := bulkgen.GenerateRandomLocations(baseLocation, size)

	var current partner.Location
	for _, l := range got {
		if reflect.DeepEqual(current, l) {
			t.Fatalf("GenerateRandomLocations() = %v, is duplicate, must be unique %v = %v", l, l, current)
		}
		current = l
	}
}

func TestGenerateRandomLocations_size(t *testing.T) {

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
