# Simple HTTP muxer with greedy params

This is a simple HTTP muxer. It works like [pat](https://github.com/bmizerany/pat) with important difference - all url params greedy, see above.

## Why another router???

I just need a simple router that can work on nested routes `/head/elem1/elem2/.../elemN/tail`. That's all!

Routers like [pat](https://github.com/bmizerany/pat), [chi](https://github.com/go-chi/chi), [echo](https://github.com/labstack/echo) are awesome, but they cannot cover this particular use case.

## Match example

<table>
<thead>
<tr><th>Pattern</th><th>Path</th><th>Match?</th><th>Params</th></tr>
</thead>
<tbody>
<tr><td>/simple</td><td>/simple/test</td><td>Yes</td><td>{}</td></tr>
<tr><td>/simple</td><td>/s1mp1e</td><td>No</td><td>{}</td></tr>
<tr><td>/one/:param1/three</td><td>/one/two/three</td><td>Yes</td><td>{param1:"two"}</td></tr>
<tr><td>/one/:param1/five</td><td>/one/two/three/four/five</td><td>Yes</td><td>{param1:"two/three/four"}</td></tr>
<tr><td>/1/:param1/5/:param2/10/:param3</td><td>/1/2/3/4/5/6/7/8/9/10/11/12</td><td>Yes</td><td>{param1:"2/3/4", param2:"6/7/8/9", param3:"11/12"}</td></tr>
</tbody>
</table>

## Params

```go
func(w http.ResponseWriter, r *http.Request) {
	params := mux.GetParams(r)
...
``` 


## Example

```go
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
```