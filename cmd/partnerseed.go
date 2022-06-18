package partnerseed

import (
	"log"

	"github.com/pysf/stunning-couscous/internal/bulkgen"
	"github.com/pysf/stunning-couscous/internal/partner"
)

func Start() {

	locations := bulkgen.GenerateRandomLocations(partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}, 100)

	partners := bulkgen.GeneratePartner(locations)

	partnerRepo, err := partner.NewPostgreRepo()
	if err != nil {
		log.Fatal(err)
	}

	partnerRepo.BulkImport(partners)

}
