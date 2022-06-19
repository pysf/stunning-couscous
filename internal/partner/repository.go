package partner

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"

	"github.com/lib/pq"
	"github.com/pysf/stunning-couscous/internal/db"
)

type Repository interface {
	BulkImport([]Partner) error
	FindBestMatch(context.Context, Location, string) ([]Partner, error)
	GetPartner(context.Context, int64) (*Partner, error)
}

func NewPartnerRepo() (Repository, error) {

	db, err := db.NewPostgreConnection()
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

	p, err := scanRow(*row)
	if err != nil {
		return nil, fmt.Errorf("GetPartner: scanRow err= %w", err)
	}

	return p, nil
}

func (ps PartnerRepo) FindBestMatch(ctx context.Context, l Location, experience string) ([]Partner, error) {

	// earthdistance is used, https://www.postgresql.org/docs/current/earthdistance.html
	rows, err := ps.DB.QueryContext(ctx, `
			SELECT
				id,
				CAST(distance AS INT),
				rating,
				experiences,
				operatingradius
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
		if err := rows.Scan(&p.ID, &p.Distance, &p.Rating, pq.Array(&p.Experiences), &p.OperatingRadius); err != nil {
			return nil, fmt.Errorf("FindBestMatch: extract row err= %w", err)
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
	ID              int64
	Rating          int
	OperatingRadius int
	Distance        int
	Experiences     []string
	Location
}

func scanRow(row sql.Row) (*Partner, error) {
	p := &Partner{}
	var point string
	err := row.Scan(&p.ID, pq.Array(&p.Experiences), &p.OperatingRadius, &p.Rating, &point)
	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		if err = p.parsePostgresPoint(point); err != nil {
			return nil, fmt.Errorf("scanRow: parse point err= %w", err)
		}
		return p, nil
	default:
		return nil, fmt.Errorf("scanRow: err= %w", err)
	}
}

type Location struct {
	Latitude  float64
	Longitude float64
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
