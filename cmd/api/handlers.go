package main

import (
	"io/ioutil"
	"net/http"
	"strings"
	"unicode/utf8"

	"inmemcache/pkg/cache"
)

// handleAPI http request handler for cache api
func (a *apiServer) handleAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")
		if len(path) != 3 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			value, err := a.cache.Get(path[2])
			if err != nil {
				if err == cache.ErrorKeyNotFound {
					w.WriteHeader(http.StatusNotFound)
					return
				} else if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					a.errorLog.Fatalf("error occurred processing the request: %w", err)
					return
				}
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(value))

		case http.MethodPost:
			requestBody, err := ioutil.ReadAll(r.Body)

			// validate if request body contains valid utf8 sequence
			valid := utf8.Valid(requestBody)
			if !valid {
				w.WriteHeader(http.StatusBadRequest)
				a.errorLog.Fatal("invalid utf8 sequence as body")
				return
			}

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				a.errorLog.Fatalf("error occurred while reading the body: %w", err)
				return
			}

			err = a.cache.Set(path[2], string(requestBody))
			if err != nil {
				w.WriteHeader(http.StatusOK)
				return
			}

		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}
