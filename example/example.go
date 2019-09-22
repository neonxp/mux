// +build ignore

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/neonxp/mux"
)

func main() {
	m := mux.New()

	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware", "hello!")
			next.ServeHTTP(w, r)
		})
	}

	m.Get("/head/:param1/middle/:param2/tail", func(writer http.ResponseWriter, reader *http.Request) {
		params := mux.GetParams(reader)
		mux.Error(mux.Plain(fmt.Sprintf("param1=%s, param2=%s", params["param1"], params["param2"]), writer), writer)
	}, middleware)

	http.Handle("/", m)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
