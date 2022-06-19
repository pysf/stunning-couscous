package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
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

	p, err := s.prtner.GetPartner(r.Context(), numId)
	if err != nil {
		return fmt.Errorf("GetPartner: err= %w", err)
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
