package main

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPractice(t *testing.T) {
	r := require.New(t)

	webApp := createAppRoutes()

	// Send 4 sequential requests: 2 to /foo and 2 to /bar
	// 1 request to /foo and 1 request to bar should be successful
	// other requests should be rejected with 429 status code
	testCases := []struct {
		path         string
		expectedCode int
	}{
		{"/foo", http.StatusOK},
		{"/bar", http.StatusOK},
		{"/foo", http.StatusTooManyRequests},
		{"/bar", http.StatusTooManyRequests},
		{"/foo", http.StatusTooManyRequests},
		{"/bar", http.StatusTooManyRequests},
	}

	for _, tc := range testCases {
		req := request(r, tc.path)
		resp, tErr := webApp.Test(req)
		r.NoError(tErr)

		r.Equal(tc.expectedCode, resp.StatusCode)
		r.True(IsValidUUID(resp.Header.Get(fiber.HeaderXRequestID)))
	}
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func request(r *require.Assertions, path string) *http.Request {
	req, tErr := http.NewRequest(http.MethodGet, path, nil)
	r.NoError(tErr)

	return req
}
