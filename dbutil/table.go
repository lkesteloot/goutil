// Copyright 2013 HeadCode

package dbutil

import (
	"fmt"
	"strings"
)

type Table struct {
	// Name of table.
	name string

	// All fields in the table except "id" and "created_at".
	fields []string

	// Comma-separated version of fields.
	partialFields string

	// Like partialFields but with "id" at the front and "created_at" at the end.
	AllFields string

	// Query string for inserting and returning the ID.
	insertQuery string

	// Query string for updating.
	updateQuery string
}

// Create a new Table with the specified name and field list. The field list must
// not include "id" or "created_at".
func MakeTable(name string, fields []string) *Table {
	table := &Table{
		name:   name,
		fields: fields,
	}

	table.partialFields = strings.Join(fields, ",")
	table.AllFields = "id," + table.partialFields + ",created_at"
	table.insertQuery = `INSERT INTO "` + table.name + `" (` + table.partialFields + `) VALUES (` +
		generateParameters(len(fields)) + `) RETURNING id`
	table.updateQuery = `UPDATE "` + table.name + `" SET ` + generateNamedParameters(fields) +
		fmt.Sprintf(` WHERE id = $%d`, len(fields)+1)

	return table
}

// Insert a new row. Include all fields except "id" and "created_at". Returns the
// new ID or the error.
func (t *Table) Insert(tx *Tx, params ...interface{}) (id IdField, err error) {
	if len(params) != len(t.fields) {
		panic(fmt.Sprintf("Wrong number of arguments (%d instead of %d)",
			len(params), len(t.fields)))
	}

	row := tx.QueryRow(t.insertQuery, params...)
	err = row.Scan(&id)
	return
}

// Updates a row. Include all fields except "id" and "created_at", followed by "id".
func (t *Table) MustUpdate(tx *Tx, params ...interface{}) {
	if len(params) != len(t.fields)+1 {
		panic(fmt.Sprintf("Wrong number of arguments (%d instead of %d)",
			len(params), len(t.fields)+1))
	}

	tx.MustExec(t.updateQuery, params...)
}

// Updates a row. Include all fields except "id" and "created_at", followed by "id".
func (t *Table) Update(tx *Tx, params ...interface{}) error {
	if len(params) != len(t.fields)+1 {
		panic(fmt.Sprintf("Wrong number of arguments (%d instead of %d)",
			len(params), len(t.fields)+1))
	}

	_, err := tx.Exec(t.updateQuery, params...)
	return err
}

// Generate a comma-separated list of $1, $2, $3, ..., $count parameters
func generateParameters(count int) string {
	var params []string

	for i := 1; i <= count; i++ {
		params = append(params, fmt.Sprintf("$%d", i))
	}

	return strings.Join(params, ",")
}

// Generate a comma-separated list of field1=$1, field2=$2, ...
func generateNamedParameters(fields []string) string {
	var params []string

	for i, field := range fields {
		params = append(params, fmt.Sprintf(`%s=$%d`, field, i+1))
	}

	return strings.Join(params, ",")
}
