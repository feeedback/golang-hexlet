// https://ru.hexlet.io/courses/go-web-development/lessons/http-standard/exercise_unit
package main

import (
	"net/http"
	"strconv"
)

var courses = map[int64]string{
	1: "Introduction to programming",
	2: "Introduction to algorithms",
	3: "Data structures",
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/courses/description", CourseDescHandler)

	http.ListenAndServe(":8080", nil)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Go to /courses/description"))
}

func CourseDescHandler(w http.ResponseWriter, r *http.Request) {
	// BEGIN (write your solution here)
	courseIDStr := r.URL.Query().Get("course_id")
	if courseIDStr == "" {
		http.Error(w, "Missing course_id parameter", http.StatusBadRequest)
		return
	}

	courseID, err := strconv.ParseInt(courseIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid course_id parameter", http.StatusBadRequest)
		return
	}

	courseDescription, ok := courses[courseID]
	if !ok {
		http.Error(w, "Course not found", http.StatusNotFound)
		return
	}

	w.Write([]byte(courseDescription))
	// END
}
