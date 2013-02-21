// Copyright 2013 HeadCode

package dbutil

import (
	"fmt"
	"strconv"
	"strings"
)

// Stores money in pennies. Could be null.
type Money struct {
	pennies int
	isNull bool
}

// For sql.Scanner interface:
func (i *Money) Scan(src interface{}) error {
	switch s := src.(type) {
	case nil:
		*i = Money{0, true}
	case string:
		intValue, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*i = Money{intValue, false}
	case int64:
		*i = Money{int(s), false}
	default:
		panic("Unknown type")
	}

	return nil
}

// Converts the value to an object that can be written to a database.
func (i Money) Print() interface{} {
	if i.isNull {
		return nil
	}

	return i.pennies
}

// Converts a string to Money. This is the inverse of ToTextField().
func ParseMoney(s string) Money {
	if s == "" {
		return Money{0, true}
	}

	// Strip dollar sign and whitespace.
	s = strings.Trim(s, " $")

	// Find decimal point.
	i := strings.Index(s, ".")
	pennies := 0
	if i >= 0 {
		var err error
		pennies, err = strconv.Atoi(s[i+1:])
		if err != nil {
			pennies = 0
		}
		s = s[:i]
	}
	dollars, err := strconv.Atoi(s)
	if err != nil {
		dollars = 0
	}

	return Money{dollars*100 + pennies, false}
}

// Return the plain text to show in an HTML text field.
func (m *Money) ToTextField() string {
	if m.isNull {
		return ""
	}

	return fmt.Sprintf("$%d.%02d", m.pennies/100, m.pennies%100)
}
