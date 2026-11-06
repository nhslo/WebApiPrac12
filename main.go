package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Student struct {
	ID          int
	Name, Group string
	GPA         float64
}

var students = []Student{

	{ID: 1, Name: "Ivan Petrov", Group: "CS-301", GPA: 3.7},

	{ID: 2, Name: "Anna Kim", Group: "CS-301", GPA: 3.9},

	{ID: 3, Name: "Maria Ivanova", Group: "CS-302", GPA: 3.2},

	{ID: 4, Name: "Daniyar Khan", Group: "CS-302", GPA: 2.8},

	{ID: 5, Name: "Sergey Sidorov", Group: "CS-303", GPA: 3.5},
}

func send(responseWriter http.ResponseWriter, minimumGPA any) {
	responseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(responseWriter).Encode(minimumGPA)
}

func home(responseWriter http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(responseWriter, request)
		return
	}
	fmt.Fprintln(responseWriter, "Student Directory: /students, /students/{id}, /students/group/{group}")
}

// При наличии minGPA оставляем студентов с подходящим баллом.
func list(responseWriter http.ResponseWriter, request *http.Request) {
	result := []Student{}
	value := request.URL.Query().Get("minGPA")
	min := 0.0
	if value != "" {
		minimumGPA, operationError := strconv.ParseFloat(value, 64)
		if operationError != nil || minimumGPA < 0 || minimumGPA > 4 {
			http.Error(responseWriter, "Invalid GPA", 400)
			return
		}
		min = minimumGPA
	}
	for _, student := range students {
		if student.GPA >= min {
			result = append(result, student)
		}
	}
	send(responseWriter, result)
}

// Получаем ID из адреса и ищем студента.
func get(responseWriter http.ResponseWriter, request *http.Request) {
	id, operationError := strconv.Atoi(request.PathValue("id"))
	if operationError != nil {
		http.Error(responseWriter, "Invalid ID", 400)
		return
	}
	for _, student := range students {
		if student.ID == id {
			send(responseWriter, student)
			return
		}
	}
	http.NotFound(responseWriter, request)
}

func group(responseWriter http.ResponseWriter, request *http.Request) {
	result := []Student{}
	for _, student := range students {
		if student.Group == request.PathValue("group") {
			result = append(result, student)
		}
	}
	send(responseWriter, result)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", home)
	mux.HandleFunc("GET /students", list)
	mux.HandleFunc("GET /students/{id}", get)
	mux.HandleFunc("GET /students/group/{group}", group)
	log.Println("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
