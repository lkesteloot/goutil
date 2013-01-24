package dbutil

import (
	"database/sql"
	"fmt"
	_ "github.com/bmizerany/pq"
	"log"
)

// Wrap this type so that we can panic() in error cases.
type Tx sql.Tx

// Like the Exec() function of sql.Tx.
func (tx *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return ((*sql.Tx)(tx)).Exec(query, args...)
}

// Like the Exec() function of sql.Tx, but calls panic() on error.
func (tx *Tx) MustExec(query string, args ...interface{}) sql.Result {
	r, err := ((*sql.Tx)(tx)).Exec(query, args...)
	if err != nil {
		// Panic so that we have the stack trace and the error message.
		panic(err)
	}

	return r
}

// Like the Query() function of sql.Tx, but calls panic() on error.
func (tx *Tx) MustQuery(query string, args ...interface{}) *sql.Rows {
	r, err := ((*sql.Tx)(tx)).Query(query, args...)
	if err != nil {
		// Panic so that we have the stack trace and the error message.
		panic(err)
	}

	return r
}

// Like the QueryRow() function of sql.Tx.
func (tx *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	return ((*sql.Tx)(tx)).QueryRow(query, args...)
}

func DoWithDb(credentials string, callback func(*sql.DB)) {
	db, err := sql.Open("postgres", credentials)
	if err != nil {
		log.Printf("Can't open database: %s", err)
		return
	}
	defer db.Close()

	callback(db)
}

func DoWithTx(credentials string, callback func(*Tx)) {
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
				err = tx.Commit()
				if err != nil {
					panic(err)
				}
			}
		}()

		callback((*Tx)(tx))
	})
}

func GetSchemaVersion(tx *Tx) (version int) {
	// First see if the schema_tracker table exists.
	var rowCount int
	row := tx.QueryRow(`
		SELECT count(*)
		FROM information_schema.tables
		WHERE table_name = 'schema_tracker'`)
	err := row.Scan(&rowCount)
	if err != nil {
		// Shouldn't happen.
		panic(err)
	}

	if rowCount == 0 {
		// Relation doesn't exist, create it.
		log.Printf("The schema_tracker relation doesn't exist. Creating it.")

		// This is a double float in seconds, showing 0 digits of sub-second
		// precision, and containing a time zone. We do want the "with time zone"
		// qualifier, because it forces the timestamp to be stored in UTC time.
		// Without it, it stores whatever value you gave it, which might be local
		// time or whatever. See sections 8.5.1.3 (TIMESTAMP) and 9.9.3 (WITH TIME
		// ZONE) of the PostgreSQL reference manual.
		tx.MustExec(`
			CREATE DOMAIN global_timestamp
			AS timestamp(0) with time zone`)

		tx.MustExec(`
			CREATE TABLE schema_tracker (
				version integer PRIMARY KEY,
				comment text NOT NULL,
				created_at global_timestamp NOT NULL DEFAULT now()
			)`)

		// 0 means none.
		return 0
	}

	// Get the version.
	row = tx.QueryRow(`
		SELECT version
		FROM schema_tracker
		ORDER BY version DESC
		LIMIT 1`)
	err = row.Scan(&version)
	if err == sql.ErrNoRows {
		// Empty table.
		return 0
	} else if err != nil {
		// Shouldn't happen.
		panic(err)
	}

	return
}

func AddVersion(tx *Tx, version int, comment string) {
	log.Printf("Applying database schema version %d (%s)", version, comment)

	tx.MustExec(`
		INSERT INTO schema_tracker (version, comment)
		VALUES ($1, $2)`, version, comment)
}
