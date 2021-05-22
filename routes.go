package main

import (
	"net/http"
	"regexp"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func NewRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method: method, regex: regexp.MustCompile("^" + pattern + "$"), handler: handler}
}

var routes = []route{
	NewRoute("GET", "/", healthCheckHandler),
	NewRoute("GET", "/healthy", healthCheckHandler),
	NewRoute("POST", "/cache-file", cacheFile),
	NewRoute("GET", "/cache-file-status", getCachedFileStatus),
	NewRoute("GET", "/get-file", getFile),
	NewRoute("GET", "/invalidate-cache", invalidateCache),
}

func extractURI(uri string) string {
	for i, c := range uri {
		if c == '?' {
			runes := []rune(uri)
			return string(runes[0:i])
		}
	}
	return uri
}