package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMiddleware_RateLimitingMiddleware(t *testing.T) {
	ctx := context.Background()

	rules := []LimitRules{
		{
			Limit:    2,
			Window:   5 * time.Second,
			IP:       []string{"127.0.0.1"},
			RuleType: IPRule,
		},
	}

	mw := New(ctx, rules)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	t.Run("RequestUnderRateLimit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "127.0.0.1"
		w := httptest.NewRecorder()
		mw.RateLimitingMiddleware(handler).ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", w.Code)
		}
	})

	t.Run("RequestExceedRateLimit", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.RemoteAddr = "127.0.0.1"
			w := httptest.NewRecorder()
			mw.RateLimitingMiddleware(handler).ServeHTTP(w, req)
			if i == 2 && w.Code != http.StatusTooManyRequests {
				t.Errorf("Expected status TooManyRequests, got %v", w.Code)
			}
		}

	})

	t.Run("RequestFromDifferentIPs", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "127.0.0.2:11111"
		w := httptest.NewRecorder()
		mw.RateLimitingMiddleware(handler).ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", w.Code)
		}
	})
}
