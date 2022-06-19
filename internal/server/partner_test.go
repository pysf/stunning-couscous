package server_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/lib/pq"
	"github.com/pysf/stunning-couscous/internal/partner"
	"github.com/pysf/stunning-couscous/internal/server"
	"github.com/pysf/stunning-couscous/internal/testutils"
)

func TestServer_GetPartner(t *testing.T) {
	db, tearDown := testutils.SetupDB(t)
	defer tearDown()

	testCases := []struct {
		want partner.Partner
		w    *httptest.ResponseRecorder
		r    *http.Request
		ps   httprouter.Params
	}{
		{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest("GET", "/api/partner", nil),
			ps: httprouter.Params{
				httprouter.Param{
					Key:   "id",
					Value: "101",
				},
			},
			want: partner.Partner{
				ID:              101,
				Rating:          4,
				OperatingRadius: 8,
				Experiences:     []string{"wood", "carpet"},
				Location: partner.Location{
					Latitude:  52.528971849007036,
					Longitude: 13.430548464498173,
				},
			},
		},
	}

	for _, tt := range testCases {

		if _, err := db.Exec(`INSERT INTO  partner ("id", "experiences", "operatingradius", "rating", "location") VALUES($1, $2, $3, $4, POINT( $5, $6 ) )`,
			tt.want.ID, pq.Array(tt.want.Experiences), tt.want.OperatingRadius, tt.want.Rating, tt.want.Latitude, tt.want.Longitude); err != nil {
			t.Errorf("failed to add test partner %s", err)
		}

		server := &server.Server{
			PartnerRepo: partner.PartnerRepo{
				DB: db,
			},
		}

		if err := server.GetPartner(tt.w, tt.r, tt.ps); err != nil {
			t.Errorf("Server.GetPartner() error = %s", err)
		}

		resp := tt.w.Result()
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Server.GetPartner() read body error = %s", err)
		}

		got := partner.Partner{}
		if jsonErr := json.Unmarshal(body, &got); jsonErr != nil {
			t.Fatalf("Server.GetPartner() unmarshal json error = %s", jsonErr)
		}

		if !reflect.DeepEqual(tt.want, got) {
			t.Fatalf("Server.GetPartner() = %v; want = %v", got, tt.want)
		}
	}
}
