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

	m.Get("/hello/(.+?)/(.+?)/world", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.Context().Value("params").([]string)
		w.Write([]byte(fmt.Sprintf("First param: %s, second param: %s. All path: %s", params[1], params[2], params[0])))
	}))

	http.Handle("/", m)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
