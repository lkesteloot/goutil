// Copyright 2013 HeadCode

package dbutil

import (
	"fmt"
)

// Stores a boolean. Could be null.
type Boolean struct {
	Value  bool
	IsNull bool
}

// Return a valid boolean with the specified value.
func NewBoolean(value bool) Boolean {
	return Boolean{
		Value:  value,
		IsNull: false,
	}
}

// Return a null boolean.
func NullBoolean() Boolean {
	return Boolean{
		Value:  false,
		IsNull: true,
	}
}

// For sql.Scanner interface:
func (i *Boolean) Scan(src interface{}) error {
	switch s := src.(type) {
	case nil:
		*i = NullBoolean()
	case bool:
		*i = NewBoolean(s)
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
		return NullBoolean()
	}

	return NewBoolean(s == "true")
}

// Converts a string to Boolean, defaulting to False.
func ParseBooleanDefaultFalse(s string) Boolean {
	if s == "" {
		return NewBoolean(false)
	}

	return NewBoolean(s == "true")
}

// Converts a string to Boolean, defaulting to True.
func ParseBooleanDefaultTrue(s string) Boolean {
	if s == "" {
		return NewBoolean(true)
	}

	return NewBoolean(s == "true")
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
