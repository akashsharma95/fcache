package main

import (
	"fmt"
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
			errMsg := fmt.Errorf("incorrect path param expected len 3 got: %d", len(path))
			a.errorLog.Println(errMsg)

			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(errMsg.Error()))

			return
		}

		switch r.Method {
		case http.MethodGet:
			value, err := a.cache.Get(path[2])
			if err != nil {
				if err == cache.ErrorKeyNotFound {
					w.WriteHeader(http.StatusNotFound)

					return
				}
				errMsg := fmt.Errorf("error occurred processing the request: %w", err)
				a.errorLog.Println(errMsg)

				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(errMsg.Error()))

				return
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(value))

		case http.MethodPost:
			requestBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				errMsg := fmt.Errorf("error occurred while reading the body: %w", err)
				a.errorLog.Println(errMsg)

				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(errMsg.Error()))

				return
			}

			// validate if request body contains valid utf8 sequence
			if !utf8.Valid(requestBody) {
				errMsg := fmt.Errorf("invalid utf8 sequence as body")
				a.errorLog.Println(errMsg)

				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(errMsg.Error()))

				return
			}

			err = a.cache.Set(path[2], string(requestBody))
			if err != nil {
				w.WriteHeader(http.StatusOK)

				return
			}

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)

			return
		}
	}
}
