package partner

import (
	"database/sql"
	"fmt"
)

const createTableQuery = `
	CREATE TABLE IF NOT EXISTS partner (
		id BIGSERIAL NOT NULL PRIMARY KEY,
		experiences varchar[],
		operatingradius int,
		rating int,
		location point NOT NULL
	);
`
const installCube = `create extension if not exists cube;`
const installEarthdistance = `create extension if not exists earthdistance;`

func ApplySchema(db *sql.DB) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	//todo: add index for id
	//todo: add index for location experience

	if _, err = tx.Exec(installCube); err != nil {
		if rbkErr := tx.Rollback(); rbkErr != nil {
			return fmt.Errorf("applySchema: failed rollback install cube %w", err)
		}
		return fmt.Errorf("applySchema: install cube err= %w", err)
	}

	if _, err = tx.Exec(installEarthdistance); err != nil {
		if rbkErr := tx.Rollback(); rbkErr != nil {
			return fmt.Errorf("applySchema: rollback install earthdistance  err= %w", err)
		}
		return fmt.Errorf("applySchema: install earthdistance err= %w", err)
	}

	if _, err = tx.Exec(createTableQuery); err != nil {
		if rbkErr := tx.Rollback(); rbkErr != nil {
			return fmt.Errorf("applySchema:  rollback err= %w", err)
		}
		return fmt.Errorf("applySchema: apply schema err= %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("applySchema: Commit err= %w", err)
	}

	return nil

}
