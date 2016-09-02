package test

import (
	"net/http"
	"strings"
)

type Requests struct {
	Method string
	UrlStr string
	Header map[string]string
	Raw    string
}

func TRollingCurl(r Requests) *http.Request {
	req, _ := http.NewRequest(
		r.Method,
		r.UrlStr,
		strings.NewReader(r.Raw),
	)

	for hkey, hval := range r.Header {
		req.Header.Set(hkey, hval)
	}

	return req
}
