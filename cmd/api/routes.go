package main

import (
	"net/http"
)

// routes register all the routes for cache api
func (a *apiServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.DefaultServeMux)
	mux.HandleFunc("/key/", a.handleAPI())

	return mux
}
