// Copyright 2013 HeadCode

package webutil

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Map from URL pattern to handler function.
type DispatchHandlerMap map[string]http.HandlerFunc

// A handler function with an additional integer parameter, usually for an entity ID.
type IntegerHandlerFunc func(http.ResponseWriter, *http.Request, int)

// Map from URL pattern to handler function with integer parameter.
type DispatchIntegerHandlerMap map[string]IntegerHandlerFunc

// Create a handler that dispatches based on a sequence of maps. Each map is
// inspected in turn, and must be one of DispatchHandlerMap or DispatchIntegerHandlerMap.
// The key of each map is a string of the form "method url", where method is "GET",
// "POST", "DELETE", etc. and url is a URL or URL pattern. A pattern can include a
// single %d if the map is a DispatchIntegerHandlerMap, and the value at the %d
// will be parsed and passed to the handler.
func DispatchHandler(maps ...interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, genericMap := range maps {
			switch m := genericMap.(type) {
			case DispatchHandlerMap:
				for pattern, function := range m {
					fields := strings.SplitN(pattern, " ", 2)
					if fields[0] == r.Method && fields[1] == r.URL.Path {
						function(w, r)
						return
					}
				}
			case DispatchIntegerHandlerMap:
				for pattern, function := range m {
					fields := strings.SplitN(pattern, " ", 2)
					urlPattern := "^" + strings.Replace(regexp.QuoteMeta(fields[1]),
						"%d", "([0-9]+)", 1) + "$"
					urlRegexp := regexp.MustCompile(urlPattern)
					if fields[0] == r.Method {
						fields = urlRegexp.FindStringSubmatch(r.URL.Path)
						if len(fields) == 2 {
							number, err := strconv.Atoi(fields[1])
							if err == nil {
								function(w, r, number)
								return
							}
						}
					}
				}
			}
		}

		log.Printf("Unknown request %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	})
}
