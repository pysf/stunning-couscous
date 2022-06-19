package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pysf/stunning-couscous/internal/partner"
)

func NewServer() (*Server, error) {
	partnerRepo, err := partner.NewPartnerRepo()
	if err != nil {
		return nil, fmt.Errorf("NewServer:  server err = %w", err)
	}

	return &Server{
		prtner: partnerRepo,
	}, nil
}

type Server struct {
	prtner partner.Repository
}

func (s *Server) Start() {

	router := httprouter.New()

	router.GET("/api/partner/:id", wrapWithErrorHandler(s.GetPartner))
	fmt.Println("Starting...")
	log.Fatal(http.ListenAndServe(":8080", router))

}
