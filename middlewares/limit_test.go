package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLimitRate(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a new request recorder
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Test with a rate limit of 1 request per second
	rl := NewRateLimiter(1, time.Second)
	limitedHandler := LimitRate(handler, rl)

	// First request should succeed
	limitedHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Second request should be rate limited
	w = httptest.NewRecorder()
	limitedHandler(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status code %d, got %d", http.StatusTooManyRequests, w.Code)
	}
}
