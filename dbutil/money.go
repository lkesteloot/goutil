// Copyright 2013 HeadCode

package dbutil

import (
	"fmt"
	"strconv"
	"strings"
)

// Stores money in pennies. Could be null.
type Money struct {
	Pennies int
	IsNull  bool
}

// Make a new money with pennies.
func NewMoney(pennies int) Money {
	return Money{
		Pennies: pennies,
		IsNull:  false,
	}
}

// Make a null money structure.
func NullMoney() Money {
	return Money{
		Pennies: 0,
		IsNull:  true,
	}
}

// For sql.Scanner interface:
func (i *Money) Scan(src interface{}) error {
	switch s := src.(type) {
	case nil:
		*i = NullMoney()
	case string:
		intValue, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*i = NewMoney(intValue)
	case int64:
		*i = NewMoney(int(s))
	default:
		panic("Unknown type")
	}

	return nil
}

// Converts the value to an object that can be written to a database.
func (i Money) Print() interface{} {
	if i.IsNull {
		return nil
	}

	return i.Pennies
}

// Converts a string to Money. This is the inverse of ToTextField().
func ParseMoney(s string) Money {
	if s == "" {
		return NullMoney()
	}

	// Strip dollar sign and whitespace.
	s = strings.Trim(s, " $")

	// Find decimal point.
	i := strings.Index(s, ".")
	pennies := 0
	if i >= 0 {
		var err error
		penniesStr := s[i+1:]
		for len(penniesStr) < 2 {
			penniesStr = penniesStr + "0"
		}
		if len(penniesStr) > 2 {
			penniesStr = penniesStr[:2]
		}
		pennies, err = strconv.Atoi(penniesStr)
		if err != nil {
			pennies = 0
		}
		s = s[:i]
	}
	dollars, err := strconv.Atoi(s)
	if err != nil {
		dollars = 0
	}

	return NewMoney(dollars*100 + pennies)
}

// Return the plain text to show in an HTML text field.
func (m Money) ToTextField() string {
	if m.IsNull {
		return ""
	}

	return fmt.Sprintf("$%d.%02d", m.Pennies/100, m.Pennies%100)
}

// Returns the money times the multiplier. Returns a Null Money if
// the object is Null.
func (m Money) MultipliedBy(x float32) Money {
	if m.IsNull {
		return m
	}

	return NewMoney(int(float32(m.Pennies)*x + 0.5))
}

// Returns a new Money adding this and other. If either is Null,
// returns the other.
func (m Money) Add(other Money) Money {
	if other.IsNull {
		return m
	}

	if m.IsNull {
		return other
	}

	return NewMoney(m.Pennies + other.Pennies)
}
