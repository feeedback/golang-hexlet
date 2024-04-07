// https://ru.hexlet.io/courses/go-web-development/lessons/logging/exercise_unit
package main

import (
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func SetLogger(l *logrus.Logger) {
	logger = l
}

func init() {
	logger = logrus.New()
	logger.SetOutput(os.Stdout)
}

func sumHandler(w http.ResponseWriter, r *http.Request) {
	xParam := r.URL.Query().Get("x")
	x, err := strconv.Atoi(xParam)
	if err != nil {
		w.Write([]byte("-1"))
		return
	}

	yParam := r.URL.Query().Get("y")
	y, err := strconv.Atoi(yParam)
	if err != nil {
		w.Write([]byte("-1"))
		return
	}

	sum := x + y

	if math.MaxInt-abs(x) < abs(y) {
		logger.WithFields(logrus.Fields{"x": x, "y": y}).Warn("Sum overflows int")

		w.Write([]byte("-1"))
		return
	}

	w.Write([]byte(strconv.Itoa(sum)))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go to /sum"))
	})

	// BEGIN (write your solution here)
	http.HandleFunc("/sum", sumHandler)
	// END

	port := "8080"
	logWithPort := logrus.WithFields(logrus.Fields{
		"port": port,
	})
	logWithPort.Info("Starting a web-server on port")
	logWithPort.Fatal(http.ListenAndServe(":"+port, nil))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
