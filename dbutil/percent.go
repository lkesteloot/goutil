// Copyright 2013 HeadCode

package dbutil

import (
	"fmt"
	"strconv"
	"strings"
)

// Stores a percentage. Note that the value is the actual percentage (0 to 100),
// not the normalized value 0 to 1. Could be null.
type Percent struct {
	Value float32
	IsNull bool
}

// For sql.Scanner interface:
func (i *Percent) Scan(src interface{}) error {
	switch s := src.(type) {
	case nil:
		*i = Percent{0, true}
	case string:
		floatValue, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return err
		}
		*i = Percent{float32(floatValue), false}
	case float32:
		*i = Percent{s, false}
	case float64:
		*i = Percent{float32(s), false}
	default:
		panic(fmt.Sprintf("Unknown type %T", src))
	}

	return nil
}

// Converts the value to an object that can be written to a database.
func (i Percent) Print() interface{} {
	if i.IsNull {
		return nil
	}

	return i.Value
}

// Converts a string to Percent. This is the inverse of ToTextField().
func ParsePercent(s string) Percent {
	if s == "" {
		return Percent{0, true}
	}

	// Strip percent sign and whitespace.
	s = strings.Trim(s, " %")

	// Parse as float.
	floatValue, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return Percent{0, true}
	}

	return Percent{float32(floatValue), false}
}

// Return the plain text to show in an HTML text field.
func (p *Percent) ToTextField() string {
	if p.IsNull {
		return ""
	}

	return fmt.Sprintf("%g%%", p.Value)
}
