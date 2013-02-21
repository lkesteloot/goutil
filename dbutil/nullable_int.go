// Copyright 2013 HeadCode

package dbutil

import (
	"strconv"
)

// Stores an integer that could be null. Null is represented as -1.
type NullableInt int

const nullNullableInt = NullableInt(-1)

// For sql.Scanner interface:
func (i *NullableInt) Scan(src interface{}) error {
	switch s := src.(type) {
	case nil:
		*i = nullNullableInt
	case string:
		intValue, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*i = NullableInt(intValue)
	case int64:
		*i = NullableInt(s)
	default:
		panic("Unknown type")
	}

	return nil
}

// Converts the ID to an object that can be written to a database.
func (i NullableInt) Print() interface{} {
	if i == nullNullableInt {
		return nil
	}

	return int(i)
}

func ParseNullableInt(s string) NullableInt {
	if s == "" {
		return nullNullableInt
	}

	intValue, err := strconv.Atoi(s)
	if err != nil {
		return nullNullableInt
	}

	return NullableInt(intValue)
}
