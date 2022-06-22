package bulkgen

import (
	"math"
	"math/rand"

	"github.com/pysf/stunning-couscous/internal/partner"
)

func GeneratePartner(locations []partner.Location) []partner.Partner {

	partners := make([]partner.Partner, len(locations))
	experiences := []string{"wood", "carpet", "tiles"}

	for i, loc := range locations {

		partners[i] = partner.Partner{
			OperatingRadius: rand.Intn(5000) + 1000,
			Rating:          rand.Intn(5) + 1,
			Experiences:     []string{experiences[0], experiences[rand.Intn(2)+1]},
			Location:        loc,
		}

	}

	return partners
}

func GenerateRandomLocations(baseLocation partner.Location, size int) []partner.Location {

	var distance int64 = 1
	result := make([]partner.Location, 0, size)

	for {
		distance = distance + 1

		for i := 1; i <= rand.Intn(6)+1; i++ {
			loc := shiftLocation(baseLocation, float64(distance), float64((60)*i))
			result = append(result, loc)
			if len(result) >= size {
				return result
			}
		}
	}

}

func shiftLocation(l partner.Location, distance, bearing float64) partner.Location {

	R := 6378.1                       // Radius of the Earth
	brng := bearing * math.Pi / 180   // Convert bearing to radian
	lat := l.Latitude * math.Pi / 180 // Current coords to radians
	lon := l.Longitude * math.Pi / 180
	// Do the math magic

	lat = math.Asin(math.Sin(lat)*math.Cos(distance/R) + math.Cos(lat)*math.Sin(distance/R)*math.Cos(brng))
	lon += math.Atan2(math.Sin(brng)*math.Sin(distance/R)*math.Cos(lat), math.Cos(distance/R)-math.Sin(lat)*math.Sin(lat))

	return partner.Location{
		Latitude:  (lat * 180 / math.Pi),
		Longitude: (lon * 180 / math.Pi),
	}
}
