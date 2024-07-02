package httpx

import (
	"net/http"
)

func GetTraceparent(r *http.Request) string {
	traceparent := r.Header.Get("Traceparent")
	return traceparent
}
