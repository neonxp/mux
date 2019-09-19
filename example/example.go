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

	m.Get("/head/:param1/middle/:param2/tail", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.GetParams(r)
		w.Write([]byte(fmt.Sprintf("param1=%s, param2=%s", params["param1"], params["param2"])))
	}))

	http.Handle("/", m)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
