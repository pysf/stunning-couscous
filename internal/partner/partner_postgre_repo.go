package partner

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pysf/stunning-couscous/internal/db"
)

func NewPostgreRepo() (*PostgreRepo, error) {

	db, err := db.NewPostgreConnection()
	if err != nil {
		return nil, fmt.Errorf("NewPostgreRepo(): failed to create psq connection %w", err)
	}

	if err = applySchema(db); err != nil {
		return nil, err
	}

	return &PostgreRepo{
		db: db,
	}, nil
}

type PostgreRepo struct {
	db *sql.DB
}

func (ps PostgreRepo) BulkImport(partners []Partner) error {
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

// func (ms PostgreRepo) CreatePartner(ctx context.Context, p Partner) error {

// 	_, err := ms.db.ExecContext(ctx, `
// 	INSERT INTO partner(location,experiences,rating,operatingradius)
// 	VALUES(POINT($1,$2),$3,$4,$5)
// 	`, p.Latitude, p.Longitude, pq.Array(p.Experiences), p.Rating, p.OperatingRadius)
// 	if err != nil {
// 		return fmt.Errorf("createPartner(): failed to insert partner %w", err)
// 	}
// 	return nil
// }
