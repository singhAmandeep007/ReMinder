package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/singhAmandeep007/ReMinder/backend/gin-server/pkg/logger"
	"github.com/stretchr/testify/assert"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func setupRateLimiterTest(limit int, window time.Duration) (*gin.Engine, *rateLimiterMiddleware) {
	log := logger.New()
	rl := NewRateLimiterMiddleware(log, limit, window).(*rateLimiterMiddleware)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(rl.RateLimiter())
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	return router, rl
}

func TestRateLimiter_AllowedRequestsWithinLimit(t *testing.T) {
	limit := 5
	window := time.Minute
	router, _ := setupRateLimiterTest(limit, window)

	for i := 0; i < limit; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:1234" // Set a fixed IP for testing
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "Request should be allowed within limit")
	}
}

func TestRateLimiter_BlockedRequestsOverLimit(t *testing.T) {
	limit := 3
	window := time.Minute
	router, _ := setupRateLimiterTest(limit, window)

	// Make requests up to the limit
	for i := 0; i < limit; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:1234"
		router.ServeHTTP(w, req)
	}

	// Make one more request that should be blocked
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Request should be blocked over limit")

	var response ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Rate limit exceeded", response.Message)
}

func TestRateLimiter_DifferentIPAddresses(t *testing.T) {
	limit := 2
	window := time.Minute
	router, _ := setupRateLimiterTest(limit, window)

	// Make requests from different IPs
	for i := 0; i < limit; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = fmt.Sprintf("127.0.0.%d:1234", i+1)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Requests from different IPs should be allowed")
	}
}

func TestRateLimiter_WindowExpiration(t *testing.T) {
	limit := 2
	window := 100 * time.Millisecond // Short window for testing
	router, rl := setupRateLimiterTest(limit, window)

	// Make requests up to the limit
	for i := 0; i < limit; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:1234"
		router.ServeHTTP(w, req)
	}

	// Make one more request that should be blocked
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Request should be blocked initially")

	// Wait for the window to expire
	time.Sleep(window + 10*time.Millisecond)

	// Now requests should be allowed again
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Request should be allowed after window expires")

	// Verify the client map was cleaned up
	rl.mu.Lock()
	defer rl.mu.Unlock()
	assert.Len(t, rl.clients["127.0.0.1"], 1, "Only one timestamp should remain after window expiration")
}

func TestRateLimiter_ConcurrentRequests(t *testing.T) {
	limit := 10
	window := time.Minute
	router, _ := setupRateLimiterTest(limit, window)

	var wg sync.WaitGroup
	results := make(chan int, limit*2)

	for i := 0; i < limit*2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "127.0.0.1:1234" // All from same IP
			router.ServeHTTP(w, req)
			results <- w.Code
		}()
	}

	wg.Wait()
	close(results)

	successCount := 0
	blockedCount := 0
	for code := range results {
		if code == http.StatusOK {
			successCount++
		} else if code == http.StatusTooManyRequests {
			blockedCount++
		}
	}

	assert.Equal(t, limit, successCount, "Should allow exactly limit requests")
	assert.Equal(t, limit, blockedCount, "Should block exactly limit requests when over limit")
}

func TestRateLimiter_ClientIPExtraction(t *testing.T) {
	limit := 1
	window := time.Minute
	router, rl := setupRateLimiterTest(limit, window)

	testCases := []struct {
		name       string
		remoteAddr string
		xForwarded string
		expectedIP string
	}{
		{"Direct IP", "192.168.1.1:1234", "", "192.168.1.1"},
		{"With X-Forwarded-For", "192.168.1.2:1234", "203.0.113.195, 70.41.3.18", "203.0.113.195"},
		{"Empty X-Forwarded-For", "192.168.1.3:1234", "", "192.168.1.3"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = tc.remoteAddr
			if tc.xForwarded != "" {
				req.Header.Set("X-Forwarded-For", tc.xForwarded)
			}
			router.ServeHTTP(w, req)

			rl.mu.Lock()
			_, exists := rl.clients[tc.expectedIP]
			rl.mu.Unlock()
			assert.True(t, exists, "Client IP should be correctly extracted and stored")
		})
	}
}

func TestRateLimiter_CleanupOldEntries(t *testing.T) {
	limit := 5
	window := 100 * time.Millisecond
	router, rl := setupRateLimiterTest(limit, window)

	// Make initial request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	router.ServeHTTP(w, req)

	// Verify entry exists
	rl.mu.Lock()
	assert.Len(t, rl.clients["127.0.0.1"], 1)
	rl.mu.Unlock()

	// Wait for window to expire plus some buffer
	time.Sleep(window + 10*time.Millisecond)

	// Make another request which should trigger cleanup
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	router.ServeHTTP(w, req)

	// Verify old timestamp was cleaned up
	rl.mu.Lock()
	assert.Len(t, rl.clients["127.0.0.1"], 1, "Only the new timestamp should remain")
	rl.mu.Unlock()
}

func TestRateLimiter_MultipleWindows(t *testing.T) {
	limit := 2
	window := 200 * time.Millisecond
	router, _ := setupRateLimiterTest(limit, window)
	ip := "127.0.0.1:1234"

	// First window - make exactly 'limit' requests
	for i := 0; i < limit; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, fmt.Sprintf("Request %d should be allowed", i+1))
	}

	// Verify next request is blocked
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = ip
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Request over limit should be blocked")

	// Wait for window to fully expire
	time.Sleep(window + 50*time.Millisecond) // Added more buffer time

	// After window expires, we should be able to make 'limit' requests again
	allowedCount := 0
	for i := 0; i < limit; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip
		router.ServeHTTP(w, req)
		if w.Code == http.StatusOK {
			allowedCount++
		}
	}
	assert.Equal(t, limit, allowedCount, "Should allow full limit after window expires")

	// Verify we get blocked again after using the new window's limit
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = ip
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Should block again after using new window's limit")
}
