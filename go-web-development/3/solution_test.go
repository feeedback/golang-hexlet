package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestPractice(t *testing.T) {
	testCases := []struct {
		name                 string
		x                    int
		y                    int
		expectedResponseBody string
		expectedCode         int
		expectedLog          string
	}{
		{
			name:                 "positive",
			x:                    11,
			y:                    55,
			expectedResponseBody: "66",
			expectedCode:         http.StatusOK,
			expectedLog:          "",
		},
		{
			name:                 "positive #2",
			x:                    124,
			y:                    77777,
			expectedResponseBody: "77901",
			expectedCode:         http.StatusOK,
			expectedLog:          "",
		},
		{
			name:                 "overflow with bigger x",
			x:                    math.MaxInt - 10,
			y:                    700,
			expectedResponseBody: "-1",
			expectedCode:         http.StatusOK,
			expectedLog:          fmt.Sprintf("level=warning msg=\"Sum overflows int\" x=%d y=700", math.MaxInt-10),
		},
		{
			name:                 "overflow with bigger x",
			x:                    500,
			y:                    math.MaxInt - 10,
			expectedResponseBody: "-1",
			expectedCode:         http.StatusOK,
			expectedLog:          fmt.Sprintf("level=warning msg=\"Sum overflows int\" x=500 y=%d", math.MaxInt-10),
		},
	}

	app := fiber.New()
	app.Get("/sum", adaptor.HTTPHandlerFunc(sumHandler))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, tErr := http.NewRequest(
				http.MethodGet,
				fmt.Sprintf("/sum?x=%d&y=%d", tc.x, tc.y),
				nil,
			)
			tr := require.New(t)
			tr.NoError(tErr)

			testLogger, logHook := test.NewNullLogger()
			SetLogger(testLogger)

			resp, tErr := app.Test(req)
			tr.NoError(tErr)
			defer resp.Body.Close()

			tr.Equal(tc.expectedCode, resp.StatusCode)
			if tc.expectedResponseBody != "" {
				body, rErr := io.ReadAll(resp.Body)
				tr.NoError(rErr)
				tr.Equal(tc.expectedResponseBody, string(body))
			}

			var logrusOutput string
			lastLogEntry := logHook.LastEntry()

			if lastLogEntry != nil {
				var err error
				logrusOutput, err = lastLogEntry.String()
				require.NoError(t, err)
			}

			tr.Contains(logrusOutput, tc.expectedLog)
		})
	}
}
