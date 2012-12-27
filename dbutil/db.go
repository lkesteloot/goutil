package dbutil

import (
	"database/sql"
	"fmt"
	"github.com/bmizerany/pq"
	"log"
)

func DoWithDb(credentials string, callback func(*sql.DB)) {
	db, err := sql.Open("postgres", credentials)
	if err != nil {
		log.Printf("Can't open database: %s", err)
		return
	}
	defer db.Close()

	callback(db)
}

func DoWithTx(credentials string, callback func(*sql.Tx)) {
	DoWithDb(credentials, func(db *sql.DB) {
		tx, err := db.Begin()
		if err != nil {
			panic(fmt.Sprintf("Can't create transaction (%s)", err))
		}
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				panic(r)
			} else {
				tx.Commit()
			}
		}()

		callback(tx)
	})
}

func GetSchemaVersion(tx *sql.Tx) (version int) {
	row := tx.QueryRow(`
		SELECT version
		FROM schema_tracker
		ORDER BY version DESC
		LIMIT 1`)
	err := row.Scan(&version)
	if err == sql.ErrNoRows {
		// Empty table.
		return 0
	}

	// See if we got an error.
	switch specificErr := err.(type) {
	case nil:
		return
	case pq.PGError:
		code := specificErr.Get('C')
		if code == "42P01" {
			// Relation doesn't exist, create it.
			log.Printf("The schema_tracker relation doesn't exist. Creating it.")

			// This is a double float in seconds, showing 0 digits of sub-second
			// precision, and containing a time zone. We do want the "with time zone"
			// qualifier, because it forces the timestamp to be stored in UTC time.
			// Without it, it stores whatever value you gave it, which might be local
			// time or whatever. See sections 8.5.1.3 (TIMESTAMP) and 9.9.3 (WITH TIME
			// ZONE) of the PostgreSQL reference manual.
			tx.Exec(`
				CREATE DOMAIN global_timestamp
				AS timestamp(0) with time zone`)

			tx.Exec(`
				CREATE TABLE schema_tracker (
					version integer PRIMARY KEY,
					comment text NOT NULL,
					created_at global_timestamp NOT NULL DEFAULT now()
				)`)
		} else {
			panic(fmt.Sprintf("Scan error with unknown code %s: %s", code, err))
		}
	}

	panic(fmt.Sprintf("Scan error with unknown type: %s", err))
}

func AddVersion(tx *sql.Tx, version int, comment string) {
	log.Printf("Applying database schema version %d (%s)", version, comment)

	tx.Exec(`
		INSERT INTO schema_tracker (version, comment)
		VALUES ($1, $2)`, version, comment)
}
