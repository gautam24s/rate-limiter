package middleware

import "net/http"

type MiddlewareIface interface {
	RateLimitingMiddleware(next http.HandlerFunc) http.HandlerFunc
}
