package partner

import (
	"database/sql"
	"fmt"
)

const createTableQuery = `
	CREATE TABLE IF NOT EXISTS partner (
		id serial PRIMARY KEY,
		experiences varchar[],
		operatingradius int,
		rating int,
		location point NOT NULL
	);
`
const installCube = `create extension if not exists cube;`
const installEarthdistance = `create extension if not exists earthdistance;`

func applySchema(db *sql.DB) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err = tx.Exec(installCube); err != nil {
		if rbkErr := tx.Rollback(); rbkErr != nil {
			return fmt.Errorf("applySchema: failed rollback install cube %w", err)
		}
		return fmt.Errorf("applySchema: failed install cube %w", err)
	}

	if _, err = tx.Exec(installEarthdistance); err != nil {
		if rbkErr := tx.Rollback(); rbkErr != nil {
			return fmt.Errorf("applySchema: failed rollback install earthdistance %w", err)
		}
		return fmt.Errorf("applySchema: failed install earthdistance %w", err)
	}

	if _, err = tx.Exec(createTableQuery); err != nil {
		if rbkErr := tx.Rollback(); rbkErr != nil {
			return fmt.Errorf("applySchema: failed, failde to rollback, %w", err)
		}
		return fmt.Errorf("applySchema: failed to apply schema, %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("applySchema: failed commit applySchema, %w", err)
	}

	return nil

}
