package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPractice(t *testing.T) {
	r := require.New(t)

	webApp := createAppRoutes()

	var jwtToken string
	resp := SignInResponse{}
	profileResp := ProfileResponse{}

	testCases := []struct {
		name         string
		method       string
		path         string
		jwtToken     string
		body         string
		wantCode     int
		responseBody interface{}
	}{
		{
			name:         "Signup - Success",
			method:       http.MethodPost,
			path:         "/signup",
			body:         `{"email":"test@test.com","password":"qwerty"}`,
			wantCode:     http.StatusOK,
			responseBody: nil,
		},
		{
			name:         "Signup - Conflict",
			method:       http.MethodPost,
			path:         "/signup",
			body:         `{"email":"test@test.com","password":"foobar"}`,
			wantCode:     http.StatusConflict,
			responseBody: nil,
		},
		{
			name:         "Signin - Unprocessable Entity",
			method:       http.MethodPost,
			path:         "/signin",
			body:         `{"email":"test2@test.com","password":"qwerty"}`,
			wantCode:     http.StatusUnprocessableEntity,
			responseBody: nil,
		},
		{
			name:         "Signin - Unprocessable Entity",
			method:       http.MethodPost,
			path:         "/signin",
			body:         `{"email":"test@test.com","password":"qwerty123"}`,
			wantCode:     http.StatusUnprocessableEntity,
			responseBody: nil,
		},
		{
			name:         "Signin - Success",
			method:       http.MethodPost,
			path:         "/signin",
			body:         `{"email":"test@test.com","password":"qwerty"}`,
			wantCode:     http.StatusOK,
			responseBody: &resp,
		},
		{
			name:         "Profile - Success",
			method:       http.MethodGet,
			path:         "/profile",
			wantCode:     http.StatusOK,
			responseBody: &profileResp,
		},
		{
			name:         "Profile - Unauthorized",
			method:       http.MethodGet,
			path:         "/profile",
			jwtToken:     "invalid",
			wantCode:     http.StatusUnauthorized,
			responseBody: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, tErr := http.NewRequest(
				tc.method,
				tc.path,
				strings.NewReader(tc.body),
			)
			tr := require.New(t)
			tr.NoError(tErr)

			req.Header.Set("Content-Type", "application/json")

			if tc.jwtToken != "" {
				req.Header.Set("Authorization", "Bearer "+tc.jwtToken)
			} else {
				req.Header.Set("Authorization", "Bearer "+jwtToken)
			}

			resp, tErr := webApp.Test(req)
			r.NoError(tErr)

			r.Equal(tc.wantCode, resp.StatusCode)

			if tc.responseBody != nil {
				bodyBytes, jErr := io.ReadAll(resp.Body)
				r.NoError(jErr)

				jErr = json.Unmarshal(bodyBytes, tc.responseBody)
				r.NoError(jErr)

				if signInResp, ok := tc.responseBody.(*SignInResponse); ok {
					jwtToken = signInResp.JWTToken
				}
			}
		})
	}

	r.Equal("test@test.com", profileResp.Email)
}
