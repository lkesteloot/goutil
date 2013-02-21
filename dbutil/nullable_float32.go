// Copyright 2013 HeadCode

package dbutil

import (
	"strconv"
)

// Stores a float32 that could be null. Null is represented as -1.
type NullableFloat32 float32

const nullNullableFloat32 = NullableFloat32(-1)

// For sql.Scanner interface:
func (i *NullableFloat32) Scan(src interface{}) error {
	switch s := src.(type) {
	case nil:
		*i = nullNullableFloat32
	case string:
		floatValue, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return err
		}
		*i = NullableFloat32(floatValue)
	case float32:
		*i = NullableFloat32(s)
	default:
		panic("Unknown type")
	}

	return nil
}

// Converts the ID to an object that can be written to a database.
func (i NullableFloat32) Print() interface{} {
	if i == nullNullableFloat32 {
		return nil
	}

	return float32(i)
}

func ParseNullableFloat32(s string) NullableFloat32 {
	if s == "" {
		return nullNullableFloat32
	}

	floatValue, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return nullNullableFloat32
	}
	return NullableFloat32(floatValue)
}
