// Copyright 2013 HeadCode

package webutil

import (
	"github.com/lkesteloot/goutil/dbutil"
	"net/http"
	"strconv"
)

func GetIntFormValue(r *http.Request, key string, defaultValue int) int {
	valueStr := r.FormValue(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		value = defaultValue
	}

	return value
}

// Return an IdField for a form value, or the defaultValue if the form does not
// have the specified field or if it cannot be parsed as an integer.
func GetIdFieldFormValue(r *http.Request, key string, defaultValue dbutil.IdField) dbutil.IdField {
	valueStr := r.FormValue(key)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return dbutil.IdField(value)
}
