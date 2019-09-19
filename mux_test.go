package mux

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name       string
		pattern    string
		test       string
		mustHit    bool
		mustParams map[string]string
	}{
		{"simple 1", "/simple", "/s1mp1e/test", false, map[string]string{}},
		{"simple 2", "/simple", "/simple/test", true, map[string]string{}},
		{"params 1", "/one/:middle/three", "/one/two/three", true, map[string]string{"middle": "two"}},
		{"params 2", "/one/:middle/four", "/one/two/three/four", true, map[string]string{"middle": "two/three"}},
		{"params 3", "/head/:param1/middle/:param2", "/head/one/two/middle/three/four", true, map[string]string{"param1": "one/two", "param2": "three/four"}},
		{"params 4", "/head/:param1/middle/:param2.html", "/head/one/two/middle/three/four.html", true, map[string]string{"param1": "one/two", "param2": "three/four"}},
		{"params 5", "/head/:param1/middle/prefix:param2.html", "/head/one/two/middle/prefixthree/four.html", true, map[string]string{"param1": "one/two", "param2": "three/four"}},
		{"params 6", "/head/:param1/middle/:param2/tail", "/head/one/two/middle/three/four/tail", true, map[string]string{"param1": "one/two", "param2": "three/four"}},
	}
	for _, test := range tests {
		test = test
		t.Run(test.name, func(t *testing.T) {
			m := New()
			params := map[string]string{}
			m.Get(test.pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				params = GetParams(r)
			}))
			u, _ := url.Parse("http://localhost" + test.test)
			r := new(httptest.ResponseRecorder)
			m.ServeHTTP(r, &http.Request{URL: u, Method: http.MethodGet})
			hit := r.Code != 404
			if test.mustHit != hit {
				t.Fatalf("Incorrect hit, expected %v actual %v", test.mustHit, hit)
			}
			if len(params) != len(test.mustParams) {
				t.Fatalf("Expected %d params, actual %d", len(test.mustParams), len(params))
			}
			for k, v := range test.mustParams {
				if params[k] != v {
					t.Fatalf("Expected param %s be %s, actual %s", k, v, params[k])
				}
			}
		})
	}
}
