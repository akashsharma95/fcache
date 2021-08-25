package main

import (
	"net/http"
)

func (a *apiServer) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/", http.DefaultServeMux)
	mux.HandleFunc("/key/", a.handleAPI())

	return mux
}
