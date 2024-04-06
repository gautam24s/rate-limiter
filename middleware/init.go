// Package middleware provides functionality for rate limiting requests based on IP addresses or endpoints.

package middleware

import (
	"context"
	"sync"
	"time"
)

// RuleType defines the type of rate limiting rule.
type RuleType int

const (
	// IPRule indicates a rate limiting rule based on IP addresses.
	IPRule RuleType = iota
	// EndpointRule indicates a rate limiting rule based on endpoints.
	EndpointRule
)

// Middleware represents a middleware structure for rate limiting.
type Middleware struct {
	Ctx      context.Context         // Context for cancellation
	Limiters map[string]*RateLimiter // Map of limiters for rate limiting based on IP addresses or endpoints
}

// CurrentCount represents the current count of requests and time of the last request for a given IP address.
type CurrentCount struct {
	LastReq time.Time // Time of the last request
	Counter int       // Count of requests
	IP      string    // IP address
}

// RateLimiter represents a rate limiter configuration.
type RateLimiter struct {
	Limit       int                      // Maximum number of requests allowed within a certain window
	Window      time.Duration            // Time window for rate limiting
	Endpoint    string                   // Endpoint for which rate limiting applies
	AccessMutex sync.Mutex               // Mutex for thread safety
	Count       map[string]*CurrentCount // Map to store counts of requests based on IP addresses
	IP          string                   // IP address for which rate limiting applies
}

// LimitRules defines dynamic rate limiting rules.
// IP should be provided only when limiting is desired for certain IP addresses.
type LimitRules struct {
	RuleType RuleType      // Type of rate limiting rule
	Limit    int           // Maximum number of requests allowed within a certain window
	Window   time.Duration // Time window for rate limiting
	Endpoint []string      // List of endpoints for which rate limiting applies
	IP       []string      // List of IP addresses for which rate limiting applies
}

// New creates a new instance of Middleware with provided rate limiting rules.
func New(ctx context.Context, rules []LimitRules) MiddlewareIface {
	limiters := make(map[string]*RateLimiter)

	for _, rule := range rules {
		switch rule.RuleType {
		case IPRule:
			for _, ip := range rule.IP {
				limiters[ip] = &RateLimiter{
					Limit:       rule.Limit,
					Window:      rule.Window,
					AccessMutex: sync.Mutex{},
					Count:       make(map[string]*CurrentCount),
					IP:          ip,
				}
			}
		case EndpointRule:
			for _, endpoint := range rule.Endpoint {
				limiters[endpoint] = &RateLimiter{
					Limit:       rule.Limit,
					Window:      rule.Window,
					Endpoint:    endpoint,
					AccessMutex: sync.Mutex{},
					Count:       make(map[string]*CurrentCount),
				}
			}
		}

	}
	m := &Middleware{
		Ctx:      ctx,
		Limiters: limiters,
	}
	go m.cleanUp()
	return m
}
