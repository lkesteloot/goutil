// Copyright 2013 HeadCode

package dbutil

import (
	"fmt"
)

// Stores a boolean. Could be null.
type Boolean struct {
	Value bool
	IsNull bool
}

// For sql.Scanner interface:
func (i *Boolean) Scan(src interface{}) error {
	switch s := src.(type) {
	case nil:
		*i = Boolean{false, true}
	case bool:
		*i = Boolean{s, false}
	default:
		panic(fmt.Sprintf("Unknown type %T", src))
	}

	return nil
}

// Converts the value to an object that can be written to a database.
func (i Boolean) Print() interface{} {
	if i.IsNull {
		return nil
	}

	return i.Value
}

// Converts a string to Boolean. This is the inverse of ToTextField().
func ParseBoolean(s string) Boolean {
	if s == "" {
		return Boolean{false, true}
	}

	return Boolean{s == "true", false}
}

// Return the plain text to show in an HTML text field.
func (p *Boolean) ToTextField() string {
	if p.IsNull {
		return ""
	} else if p.Value {
		return "true"
	}

	return "false"
}
