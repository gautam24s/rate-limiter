package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// getClientIP extracts the client IP address from the HTTP request.
func (m *Middleware) getClientIP(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = r.Header.Get("X-Real-Ip")
	}
	if clientIP != "" {
		return clientIP
	}
	return r.RemoteAddr
}

// RateLimitingMiddleware is an HTTP middleware function that performs rate limiting based on the configured rules.
func (m *Middleware) RateLimitingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := m.getClientIP(r)
		limiter, ok := m.Limiters[ip]
		if !ok {
			limiter, ok = m.Limiters[r.URL.Path]
		}
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		limiter.AccessMutex.Lock()
		defer limiter.AccessMutex.Unlock()

		count, ok := limiter.Count[ip]

		if !ok {
			limiter.Count[ip] = &CurrentCount{
				LastReq: time.Now(),
				Counter: 0,
				IP:      ip,
			}
			count = limiter.Count[ip]
		}
		fmt.Printf("ip: %s | count: %d | lastReq: %v \n", ip, count.Counter, count.LastReq)

		if ok && (time.Since(count.LastReq) > limiter.Window) {
			count.Counter = 0
			count.LastReq = time.Now()
		}

		if count.Counter >= limiter.Limit {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}
		count.Counter++
		next.ServeHTTP(w, r)

	}
}

// cleanUp periodically scans the rate limit records and removes outdated entries.
func (m *Middleware) cleanUp() {
	cleanupInterval := 2 * time.Minute
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-m.Ctx.Done():
			return
		case <-ticker.C:
			m.cleanupExpiredRecords()
		}
	}
}

// cleanupExpiredRecords removes outdated rate limit records.
func (m *Middleware) cleanupExpiredRecords() {
	for _, limiter := range m.Limiters {
		limiter.AccessMutex.Lock()
		for ip, count := range limiter.Count {
			if time.Since(count.LastReq) > limiter.Window {
				delete(limiter.Count, ip)
			}
		}
		limiter.AccessMutex.Unlock()
	}
}
