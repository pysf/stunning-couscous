package bulkgen

import (
	"math"
	"testing"

	"github.com/pysf/stunning-couscous/internal/partner"
)

func TestShiftLocation(t *testing.T) {

	baseLocation := partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}
	shift := 10
	got := shiftLocation(baseLocation, float64(shift), float64(60))

	distance := getDistance(got.Latitude, got.Longitude, baseLocation.Latitude, baseLocation.Longitude)
	if math.Floor(distance) != float64(shift) {
		t.Fatalf("shiftLocation() = %v, want= %v", distance, shift)
	}

}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Distance function returns the getDistance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc getDistance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// getDistance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
func getDistance(lat1, lon1, lat2, lon2 float64) float64 {

	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h)) / 1000
}
