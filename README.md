# Simple regexp based HTTP muxer

This is a simple HTTP muxer. All path patterns must be a valid go regexp (https://golang.org/pkg/regexp/syntax/).

## Why another router???

I just need a simple router that can work on nested routes `/head/elem1/elem2/.../elemN/tail`. That's all!

Routers like [pat](https://github.com/bmizerany/pat), [chi](https://github.com/go-chi/chi), [echo](https://github.com/labstack/echo) are awesome, but they cannot cover this particular use case.

## Params

If one or more groups is present in the pattern (i.e. `/book/([a-z]+)/(\d+)`) you can get the corresponding parameters from the context:

```go
func(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value("params").([]string)
...
``` 
Here `params[0]` contains substring that matches to all regexp, `params[1]` is the first group, `params[N]` is the N group.
Thus,  for regexp `/book/([a-z]+)/(\d+)` and path `/shop/book/golang/123`, the parameters will contains: `[0] => "/book/golang/123", [1] => "golang", [2] => "123"`

## Match example

<table>
<thead>
<tr><th>Pattern</th><th>Path</th><th>Params</th></tr>
</thead>
<tbody>
<tr><td>^/simple</td><td>/simple/test</td><td>[0]=>"/simple"</td></tr>
<tr><td>^/one/(.+?)/three$</td><td>/one/two/three</td><td>[0]=>"/one/two/three", [1]=>"two"</td></tr>
<tr><td>^/one/(.+?)/four$</td><td>/one/two/three/four</td><td>[0]=>"/one/two/three/four", [1]=>"two/three"</td></tr>

</tbody>
</table>

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
```