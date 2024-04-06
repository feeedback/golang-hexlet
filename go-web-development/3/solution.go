// https://ru.hexlet.io/courses/go-web-development/lessons/logging/exercise_unit
package main

import (
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sirupsen/logrus"
)

func main() {

	cwd, _ := os.Getwd()
	logFile := filepath.Join(cwd, ".log")
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()
	logger.SetOutput(file)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go to /sum"))
	})

	http.HandleFunc("/sum", func(w http.ResponseWriter, r *http.Request) {
		// BEGIN (write your solution here)
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
		// END
	})

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
