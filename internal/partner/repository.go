package partner

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pysf/stunning-couscous/internal/db"
)

type Repository interface {
	BulkImport([]Partner) error
	FindBestMatch(ctx context.Context, l Location, experience string) ([]Partner, error)
}

type Partner struct {
	ID              int
	Rating          int
	OperatingRadius int
	Distance        int
	Experiences     []string
	Location
}

type Location struct {
	Latitude  float64
	Longitude float64
}

func NewPartnerRepo() (Repository, error) {

	db, err := db.NewPostgreConnection()
	if err != nil {
		return nil, fmt.Errorf("NewPostgreRepo(): failed to create psq connection %w", err)
	}

	if err = applySchema(db); err != nil {
		return nil, err
	}

	return &partnerRepo{
		db: db,
	}, nil
}

type partnerRepo struct {
	db *sql.DB
}

func (ps partnerRepo) FindBestMatch(ctx context.Context, l Location, experience string) ([]Partner, error) {

	// earthdistance is used, https://www.postgresql.org/docs/current/earthdistance.html
	rows, err := ps.db.QueryContext(ctx, `
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
		return nil, fmt.Errorf("FindBestMatch() query failed, %w", err)
	}

	partners := make([]Partner, 0)
	for rows.Next() {
		var p Partner
		if err := rows.Scan(&p.ID, &p.Distance, &p.Rating, pq.Array(&p.Experiences), &p.OperatingRadius); err != nil {
			return nil, fmt.Errorf("FindBestMatch() extract row failed, %w", err)
		}
		partners = append(partners, p)
	}

	return partners, nil

}

func (ps partnerRepo) BulkImport(partners []Partner) error {
	txn, err := ps.db.Begin()
	if err != nil {
		return fmt.Errorf("BulkImport(): failed to begin tx, %w", err)
	}

	stmt, err := txn.Prepare(pq.CopyIn("partner", "location", "experiences", "operatingradius", "rating"))
	if err != nil {
		return fmt.Errorf("BulkImport(): failed to prepare tx, %w", err)
	}

	for _, p := range partners {
		if _, err := stmt.Exec(fmt.Sprintf("(%v,%v)", p.Latitude, p.Longitude), pq.Array(p.Experiences), p.OperatingRadius, p.Rating); err != nil {
			return fmt.Errorf("BulkImport(): exec, %w", err)
		}
	}

	if _, err := stmt.Exec(); err != nil {
		return fmt.Errorf("BulkImport(): exec finalize, %w", err)
	}

	if err := stmt.Close(); err != nil {
		return fmt.Errorf("BulkImport(): exec Close, %w", err)
	}

	if err := txn.Commit(); err != nil {
		return fmt.Errorf("BulkImport(): txn Commit, %w", err)
	}

	return nil

}
