package partner

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/lib/pq"
	"github.com/pysf/stunning-couscous/internal/db"
)

type Repository interface {
	BulkImport([]Partner) error
	FindBestMatch(ctx context.Context, loc Location, experience string) ([]Partner, error)
	GetPartner(context.Context, int64) (*Partner, error)
}

func NewPartnerRepo() (Repository, error) {

	db, err := db.NewPostgreConnection(os.Getenv("POSTGRESQL_DATABASE"))
	if err != nil {
		return nil, fmt.Errorf("NewPostgreRepo: create psq connection err= %w", err)
	}

	if err = ApplySchema(db); err != nil {
		return nil, err
	}

	return &PartnerRepo{
		DB: db,
	}, nil
}

type PartnerRepo struct {
	DB *sql.DB
}

func (ps PartnerRepo) GetPartner(ctx context.Context, id int64) (*Partner, error) {
	row := ps.DB.QueryRow(`SELECT * FROM partner WHERE ID=$1`, id)

	p := Partner{}
	if err := scanRow(*row, &p); err != nil {
		return nil, fmt.Errorf("GetPartner: scanRow err= %w", err)
	}

	return &p, nil
}

func (ps PartnerRepo) FindBestMatch(ctx context.Context, l Location, experience string) ([]Partner, error) {

	// earthdistance is used, https://www.postgresql.org/docs/current/earthdistance.html
	rows, err := ps.DB.QueryContext(ctx, `
			SELECT
				id,
				CAST(distance AS INT),
				rating,
				experiences,
				operatingradius,
				location
			FROM (
				SELECT
					id,
					location,
					operatingradius,
					experiences,
					rating,
					earth_distance (
						ll_to_earth ("location" [0],
							"location" [1]),
						ll_to_earth ($1,
							$2)) AS distance
				FROM (
					SELECT
						*
					FROM
						partner	
					WHERE
						$3 = ANY (experiences)
					) skillFiltered
				) distanceCalcusated
			WHERE
				operatingradius >= distance
			ORDER BY
				rating DESC, distance ASC; 
			`,
		l.Latitude, l.Longitude, experience)

	if err != nil {
		return nil, fmt.Errorf("FindBestMatch: run query err= %w", err)
	}
	defer rows.Close()

	partners := make([]Partner, 0)
	for rows.Next() {
		var p Partner
		var point string
		if err := rows.Scan(&p.ID, &p.Distance, &p.Rating, pq.Array(&p.Experiences), &p.OperatingRadius, &point); err != nil {
			return nil, fmt.Errorf("FindBestMatch: extract row err= %w", err)
		}

		if err = p.parsePostgresPoint(point); err != nil {
			return nil, fmt.Errorf("scanRow: parse point err= %w", err)
		}

		partners = append(partners, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("FindBestMatch: extract rows err= %w", err)
	}

	return partners, nil

}

func (ps PartnerRepo) BulkImport(partners []Partner) error {
	txn, err := ps.DB.Begin()
	if err != nil {
		return fmt.Errorf("BulkImport: begin tx err= %w", err)
	}

	stmt, err := txn.Prepare(pq.CopyIn("partner", "location", "experiences", "operatingradius", "rating"))
	if err != nil {
		return fmt.Errorf("BulkImport: prepare tx err= %w", err)
	}

	for _, p := range partners {
		if _, err := stmt.Exec(fmt.Sprintf("(%v,%v)", p.Latitude, p.Longitude), pq.Array(p.Experiences), p.OperatingRadius, p.Rating); err != nil {
			return fmt.Errorf("BulkImport: exec err= %w", err)
		}
	}

	if _, err := stmt.Exec(); err != nil {
		return fmt.Errorf("BulkImport: exec finalize err= %w", err)
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("BulkImport: exec Close err= %w", err)
	}

	if err := txn.Commit(); err != nil {
		return fmt.Errorf("BulkImport: txn Commit err= %w", err)
	}

	return nil

}

type Partner struct {
	ID              int64    `json:"id"`
	Rating          int      `json:"rating"`
	OperatingRadius int      `json:"operatingRadius"`
	Distance        int      `json:"distance"`
	Experiences     []string `json:"experiences"`
	Location        `json:"location"`
}

func scanRow(row sql.Row, p *Partner) error {

	var point string
	err := row.Scan(&p.ID, pq.Array(&p.Experiences), &p.OperatingRadius, &p.Rating, &point)
	switch err {
	case sql.ErrNoRows:
		return nil
	case nil:
		if err = p.parsePostgresPoint(point); err != nil {
			return fmt.Errorf("scanRow: parse point err= %w", err)
		}
		return nil
	default:
		return fmt.Errorf("scanRow: err= %w", err)
	}
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func FillLocation(lat, lng string, l *Location) error {
	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return fmt.Errorf("FillLocation: latitude err=%w", err)
	}

	longitude, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		return fmt.Errorf("FillLocation: longitude err=%w", err)
	}

	l.Latitude = latitude
	l.Longitude = longitude

	return nil
}

func (l *Location) parsePostgresPoint(point string) error {

	rgx := regexp.MustCompile(`([\d\.]+),([\d\.]+)`)
	res := rgx.FindStringSubmatch(point)

	if len(res) != 3 {
		return nil
	}

	lat, err := strconv.ParseFloat(res[1], 64)
	if err != nil {
		return fmt.Errorf("parsePostgresPoint: parse latitude err= %w ", err)
	}

	lng, err := strconv.ParseFloat(res[2], 64)
	if err != nil {
		return fmt.Errorf("parsePostgresPoint: parse longitude err= %w ", err)
	}

	l.Latitude = lat
	l.Longitude = lng
	return nil
}
