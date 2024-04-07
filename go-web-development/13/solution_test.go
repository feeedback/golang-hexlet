package main

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestPractice(t *testing.T) {
	t.Run("bad_request", func(t *testing.T) {
		tr := require.New(t)

		testRequest(tr, "{", fiber.StatusBadRequest, "Invalid JSON")
	})

	t.Run("panics", func(t *testing.T) {
		tr := require.New(t)

		testRequest(tr, `{"user_id":1,"message":"Hello world!"}`, fiber.StatusOK, "OK")
		testRequest(tr, `{"user_id":2,"message":"Hello world2!"}`, fiber.StatusOK, "OK")
		testRequest(tr, `{"user_id":3,"message":"Hello world3!"}`, fiber.StatusOK, "OK")
		testRequest(tr, `{"user_id":4,"message":"Hello world4!"}`, fiber.StatusInternalServerError, "Queue is full")
	})
}

func testRequest(r *require.Assertions, body string, wantCode int, wantBody string) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req, err := http.NewRequest(http.MethodPost, "/push/send", bodyReader)
	r.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	webApp := createAppRoutes()
	resp, tErr := webApp.Test(req)
	r.NoError(tErr)

	bodyBytes, err := io.ReadAll(resp.Body)
	r.NoError(err)

	r.Equal(wantCode, resp.StatusCode)
	r.Equal(wantBody, string(bodyBytes))
}
