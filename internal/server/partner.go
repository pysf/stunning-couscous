package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/pysf/stunning-couscous/internal/partner"
)

func (s *Server) GetPartner(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {
	id := ps.ByName("id")
	if len(id) == 0 {
		return NewHttpError(nil, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	numId, err := strconv.ParseInt(id, 0, 64)
	if err != nil {
		return NewHttpError(err, "invalid id format", http.StatusBadRequest)
	}

	p, err := s.PartnerRepo.GetPartner(r.Context(), numId)
	if err != nil {
		return fmt.Errorf("GetPartner: err= %w", err)
	}

	if p == nil {
		return NewHttpError(err, "partner not found", http.StatusNotFound)
	}

	res, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("GetPartner: respons to json err= %w", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("API_Version", "1.0")
	w.WriteHeader(200)
	w.Write(res)
	return nil
}

func (s *Server) FindBestMatch(w http.ResponseWriter, r *http.Request, ps httprouter.Params) error {

	bestMatchReq := BestMatchRequest{}
	FillBestMatch(r.URL.Query(), &bestMatchReq)

	if err := validator.New().Struct(bestMatchReq); err != nil {
		return NewHttpError(nil, err.Error(), http.StatusBadRequest)
	}

	l := partner.Location{}
	if err := partner.FillLocation(bestMatchReq.Latitude, bestMatchReq.Longitude, &l); err != nil {
		return fmt.Errorf("FillLocation: err= %w ", err)
	}

	// Query pagination is not implemented because it was not asked, but can be add easily
	partners, err := s.PartnerRepo.FindBestMatch(r.Context(), l, bestMatchReq.Material)
	if err != nil {
		return fmt.Errorf("PartnerRepo.FindBestMatch: err=%w", err)
	}

	body, err := json.Marshal(partners)
	if err != nil {
		return fmt.Errorf("partner list to json err=%w", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Add("API_Version", "1.0")
	w.WriteHeader(200)
	w.Write(body)

	return nil

}

type BestMatchRequest struct {
	Material  string `json:"material" validate:"required,alpha,lowercase"`
	Latitude  string `json:"latitude" validate:"required,latitude"`
	Longitude string `json:"longitude" validate:"required,longitude"`
	Square    string `json:"square" validate:"omitempty,number"`
	Phone     string `json:"phone" validate:"omitempty,e164"`
}

func FillBestMatch(url url.Values, bestMatchRequest *BestMatchRequest) {
	bestMatchRequest.Material = url.Get("material")
	bestMatchRequest.Latitude = url.Get("latitude")
	bestMatchRequest.Longitude = url.Get("longitude")
	bestMatchRequest.Square = url.Get("square")
	if len(url.Get("phone")) != 0 {
		bestMatchRequest.Phone = fmt.Sprintf("+%v", url.Get("phone"))
	}
}
