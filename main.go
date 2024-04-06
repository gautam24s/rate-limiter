// Package main demonstrates a simple HTTP server with rate limiting middleware.
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gautam24s/rate-limiter/middleware"
)

func main() {
	// Create a background context and a cancel function to gracefully stop the server.
	mainCtx := context.Background()
	ctx, cancelFunc := context.WithCancel(mainCtx)
	defer cancelFunc()

	// Define rate limiting rules for the middleware.
	limiters := []middleware.LimitRules{
		{
			Limit:    5,
			Window:   10 * time.Second,
			IP:       []string{"0.0.0.0"},
			RuleType: middleware.IPRule,
		},
		{
			Limit:    10,
			Window:   10 * time.Second,
			Endpoint: []string{"/"},
			RuleType: middleware.EndpointRule,
		},
	}

	// Create a new instance of middleware with the defined rate limiting rules.
	mw := middleware.New(ctx, limiters)

	// Define a handler function to process incoming HTTP requests.
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Request processed for %s", r.URL.Path)
	}

	// Attach the rate limiting middleware to the root endpoint.
	http.HandleFunc("/", mw.RateLimitingMiddleware(handler))

	// Start the HTTP server and listen for incoming requests.
	http.ListenAndServe(":8080", nil)
}
