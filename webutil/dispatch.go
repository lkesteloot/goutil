// Copyright 2013 HeadCode

package webutil

import (
	"fmt"
	"github.com/lkesteloot/goutil/dbutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Map from URL pattern to handler function.
type DispatchHandlerMap map[string]http.HandlerFunc

// A handler function with an additional integer parameter, usually for an entity ID.
type IdFieldHandlerFunc func(http.ResponseWriter, *http.Request, dbutil.IdField)

// Map from URL pattern to handler function with integer parameter.
type DispatchIdFieldHandlerMap map[string]IdFieldHandlerFunc

// Create a handler that dispatches based on a set of maps. Each map must be
// one of DispatchHandlerMap or DispatchIdFieldHandlerMap.
//
// The key of each map is a string of the form "method url", where method is "GET",
// "POST", "DELETE", etc. and url is a URL or URL pattern. A pattern can include a
// single %d if the map is a DispatchIdFieldHandlerMap, and the value at the %d
// will be parsed and passed to the handler.
func DispatchHandler(maps ...interface{}) http.Handler {
	// Data structures for storing maps.
	type Handler struct {
		method   string
		url      string
		function http.HandlerFunc
	}
	type IdFieldHandler struct {
		method    string
		urlRegexp *regexp.Regexp
		function  IdFieldHandlerFunc
	}
	var handlers []*Handler
	var integerHandlers []*IdFieldHandler

	// Pre-process handlers.
	for _, genericMap := range maps {
		switch m := genericMap.(type) {
		case DispatchHandlerMap:
			for pattern, function := range m {
				fields := strings.SplitN(pattern, " ", 2)
				handlers = append(handlers, &Handler{
					method:   fields[0],
					url:      fields[1],
					function: function,
				})
			}
		case DispatchIdFieldHandlerMap:
			for pattern, function := range m {
				fields := strings.SplitN(pattern, " ", 2)
				if len(fields) != 2 {
					panic("IdField handler pattern must have two fields.")
				}
				urlPattern := "^" + strings.Replace(regexp.QuoteMeta(fields[1]),
					"%d", "([0-9]+)", 1) + "$"
				urlRegexp := regexp.MustCompile(urlPattern)
				integerHandlers = append(integerHandlers, &IdFieldHandler{
					method:    fields[0],
					urlRegexp: urlRegexp,
					function:  function,
				})
			}
		default:
			panic(fmt.Sprintf("Unknown map type %T", genericMap))
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			if r.Method == handler.method && r.URL.Path == handler.url {
				handler.function(w, r)
				return
			}
		}

		for _, handler := range integerHandlers {
			if r.Method == handler.method {
				matches := handler.urlRegexp.FindStringSubmatch(r.URL.Path)
				if len(matches) == 2 {
					number, err := strconv.Atoi(matches[1])
					if err == nil {
						handler.function(w, r, dbutil.IdField(number))
						return
					}
				}
			}
		}

		log.Printf("Unknown request %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	})
}
