package httpx

import "net/http"

func GetTraceparent(r *http.Request) string {
	return r.Header.Get("Traceparent")
}
