// Copyright 2013 HeadCode

package webutil

import (
	"github.com/lkesteloot/goutil/dbutil"
	"net/http"
	"strconv"
	"strings"
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

// Return an array of IdField for a form value. The value must be
// comma-separated. Any ID that cannot be parsed is skipped.
func GetIdsFieldFormValue(r *http.Request, key string) []dbutil.IdField {
	valuesStr := strings.Split(r.FormValue(key), ",")
	ids := make([]dbutil.IdField, 0)

	for _, valueStr := range valuesStr {
		id, err := strconv.Atoi(valueStr)
		if err == nil {
			ids = append(ids, dbutil.IdField(id))
		}
	}

	return ids
}
