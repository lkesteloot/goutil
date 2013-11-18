// Copyright 2013 HeadCode

package dbutil

import (
	"strconv"
)

// Stores a database ID (either primary key or foreign key).
type IdField int

// For sql.Scanner interface:
func (id *IdField) Scan(src interface{}) error {
	switch s := src.(type) {
	case nil:
		// We use 0 to mean NULL.
		*id = 0
	case string:
		intId, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*id = IdField(intId)
	case int64:
		*id = IdField(s)
	default:
		panic("Unknown type")
	}

	return nil
}

// Converts the ID to an object that can be written to a database.
func (id IdField) Print() interface{} {
	if id == 0 {
		return nil
	}

	return int(id)
}

// Parse a field returned from a web form. If the ID can't be parsed, returns
// the default value.
func ParseIdField(s string, defaultValue IdField) IdField {
	id, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return IdField(id)
}
