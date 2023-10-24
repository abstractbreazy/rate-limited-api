package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abstractbreazy/rate-limited-api/internal/repository/redis"
	"github.com/abstractbreazy/rate-limited-api/pkg/limiter"
	redis_db "github.com/abstractbreazy/rate-limited-api/pkg/redis"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testIP = "123.45.67.57"
	tests  = []int{
		http.StatusOK,
		http.StatusOK,
		http.StatusOK,
		http.StatusOK,
		http.StatusOK,
		http.StatusTooManyRequests,
		http.StatusTooManyRequests,
		http.StatusTooManyRequests,
	}
	testLimit    uint64 = 5
	testCooldown        = 10 * time.Second
	testSubnet          = "123.45.67.0/24"
)

func prepareHandler(t *testing.T) *Handler {
	// Create a mock Redis instance
	rd, err := redis.New(&redis_db.Config{URL: "redis://localhost:6379"})
	require.NoError(t, err)

	// Create a mock limiter configuration
	limiterConfig := &limiter.Config{
		Limit:      testLimit,
		Timeout:    testCooldown,
		SubnetMask: "24",
	}

	// Create a new handler instance
	handler := NewHandler(limiterConfig, rd)
	return handler
}

func TestPingHandler(t *testing.T) {
	handler := prepareHandler(t)
	e := echo.New()
	for i := range tests {
		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add(echo.HeaderXRealIP, testIP)

		// Create a new HTTP response recorder
		rec := httptest.NewRecorder()

		// Create an Echo context
		c := e.NewContext(req, rec)

		// Call the handler function
		err := handler.Ping(c)
		assert.NoError(t, err)

		// Assert the response status code
		assert.Equal(t, tests[i], rec.Code)
	}

	// Wait for cooldown
	time.Sleep(testCooldown)

	for i := range tests {
		// Create a new HTTP request
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add(echo.HeaderXRealIP, testIP)

		// Create a new HTTP response recorder
		rec := httptest.NewRecorder()

		// Create an Echo context
		c := e.NewContext(req, rec)

		// Call the handler function
		err := handler.Ping(c)
		assert.NoError(t, err)

		// Assert the response status code
		assert.Equal(t, tests[i], rec.Code)
	}
}

func TestResetHandler(t *testing.T) {
	// Create a new Echo instance
	e := echo.New()

	handler := prepareHandler(t)

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodGet, "/reset", nil)
	req.Header.Add(echo.HeaderXRealIP, testIP)

	// Create a new HTTP response recorder
	rec := httptest.NewRecorder()

	// Create an Echo context
	c := e.NewContext(req, rec)

	// Call the handler function
	err := handler.Reset(c)

	// Assert that no error occurred
	assert.NoError(t, err)

	// Assert the response status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Assert the response body
	expectedBody := "Reset rate limit for subnet : " + testSubnet + "\n"
	assert.Equal(t, expectedBody, rec.Body.String())
}
