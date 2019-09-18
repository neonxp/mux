package mux

import (
	"net/http"
	"net/url"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name       string
		pattern    string
		test       string
		mustHit    bool
		mustParams []string
	}{
		{"simple 1", "^/simple", "/simple/test", true, []string{"/simple"}},
		{"simple 2", "^/simple", "/s1mp1e/test", false, []string{}},
		{"simple 3", "^/simple$", "/simple/test", false, []string{}},
		{"params 1", "^/one/(.+?)/three$", "/one/two/three", true, []string{"/one/two/three", "two"}},
		{"params 2", "^/one/(.+?)/four$", "/one/two/three/four", true, []string{"/one/two/three/four", "two/three"}},
	}
	for _, test := range tests {
		test = test
		t.Run(test.name, func(t *testing.T) {
			m := New()
			hit := false
			params := []string{}
			m.Get(test.pattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hit = true
				params = r.Context().Value("params").([]string)
			}))
			u, _ := url.Parse("http://localhost" + test.test)
			m.ServeHTTP(nil, &http.Request{URL: u, Method: http.MethodGet})
			if test.mustHit != hit {
				t.Fatalf("Incorrect hit, expected %v actual %v", test.mustHit, hit)
			}
			if len(params) != len(test.mustParams) {
				t.Fatalf("Expected %d params, actual %d", len(test.mustParams), len(params))
			}
			for k, v := range test.mustParams {
				if params[k] != v {
					t.Fatalf("Expected param #%d be %s, actual %s", k, v, params[k])
				}
			}
		})
	}
}
