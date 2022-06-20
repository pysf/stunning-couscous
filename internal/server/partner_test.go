package server_test

import (
	"encoding/json"
	"fmt"
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

func TestServer_FindBestMatch(t *testing.T) {
	db, tearDown := testutils.SetupDB(t)
	defer tearDown()

	baseLoc := partner.Location{
		Latitude:  52.51999140,
		Longitude: 13.40497255,
	}
	testutils.SeedTestPartners(t, db, baseLoc, 100)

	tests := []struct {
		name        string
		PartnerRepo partner.Repository
		w           *httptest.ResponseRecorder
		r           *http.Request
		ps          httprouter.Params
		query       map[string]string
		wantErr     bool
	}{
		{
			name: "Success",
			PartnerRepo: partner.PartnerRepo{
				DB: db,
			},
			w:  httptest.NewRecorder(),
			r:  httptest.NewRequest("GET", "/api/search/partner/best-match", nil),
			ps: httprouter.Params{},
			query: map[string]string{
				"material":  "wood",
				"latitude":  fmt.Sprintf("%f", baseLoc.Latitude),
				"longitude": fmt.Sprintf("%f", baseLoc.Longitude),
				"phone":     "1123456789",
				"square":    "10",
			},
			wantErr: false,
		},
		{
			name: "Validate material required",
			PartnerRepo: partner.PartnerRepo{
				DB: db,
			},
			w:  httptest.NewRecorder(),
			r:  httptest.NewRequest("GET", "/api/search/partner/best-match", nil),
			ps: httprouter.Params{},
			query: map[string]string{
				"latitude":  fmt.Sprintf("%f", baseLoc.Latitude),
				"longitude": fmt.Sprintf("%f", baseLoc.Longitude),
				"phone":     "1123456789",
				"square":    "10",
			},
			wantErr: true,
		},
		{
			name: "Validate latitude required",
			PartnerRepo: partner.PartnerRepo{
				DB: db,
			},
			w:  httptest.NewRecorder(),
			r:  httptest.NewRequest("GET", "/api/search/partner/best-match", nil),
			ps: httprouter.Params{},
			query: map[string]string{
				"material":  "wood",
				"longitude": fmt.Sprintf("%f", baseLoc.Longitude),
				"phone":     "1123456789",
				"square":    "10",
			},
			wantErr: true,
		},
		{
			name: "Validate longitude required",
			PartnerRepo: partner.PartnerRepo{
				DB: db,
			},
			w:  httptest.NewRecorder(),
			r:  httptest.NewRequest("GET", "/api/search/partner/best-match", nil),
			ps: httprouter.Params{},
			query: map[string]string{
				"material": "wood",
				"latitude": fmt.Sprintf("%f", baseLoc.Latitude),
				"phone":    "1123456789",
				"square":   "10",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &server.Server{
				PartnerRepo: tt.PartnerRepo,
			}

			q := tt.r.URL.Query()
			for k, v := range tt.query {
				q.Add(k, v)
			}
			tt.r.URL.RawQuery = q.Encode()

			if err := s.FindBestMatch(tt.w, tt.r, tt.ps); (err != nil) != tt.wantErr {
				t.Errorf("Server.FindBestMatch() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
