package main

import (
	"net/http"
	"testing"
)

func TestGetTarget(t *testing.T) {
	testCases := map[string]string{
		"http://default.default.svc:80": "domain.cai.dev",
		"http://name.namespace.type:80": "name.type.namespace.domain.cai.dev",
		"http://name.namespace.type:66": "66.name.type.namespace.clustername.cai.dev",
		"http://name.default.svc:66":    "66.name.domain.cai.dev",
	}

	for answer, item := range testCases {
		d := Destination{}
		request := http.Request{}
		request.Host = item
		d.getTarget(&request)
		if answer != d.path {
			t.Errorf("Wrong %s and %s", answer, d.path)
		}
	}

}

func BenchmarkGetTarget(t *testing.B) {

	testCases := map[string]string{
		"http://a.b.c:6": "6.a.b.c.domain.cai.dev",
	}
	d := Destination{}
	t.ResetTimer()
	for answer, item := range testCases {
		request := http.Request{}
		request.Host = item
		d.getTarget(&request)
		if answer != d.path {
			t.Errorf("Wrong %s and %s", answer, d.path)

		}
	}
}
