package main

import (
	"fmt"
	"log"

	"github.com/pysf/stunning-couscous/internal/bulkgen"
	"github.com/pysf/stunning-couscous/internal/partner"
)

func main() {

	locations := bulkgen.GenerateRandomLocations(partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}, 10)

	partners := bulkgen.GeneratePartner(locations)

	partnerRepo, err := partner.NewPartnerRepo()
	if err != nil {
		log.Fatal(err)
	}

	partnerRepo.BulkImport(partners)
	fmt.Println("hello.. to you!")

}
