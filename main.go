package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Declaring the structure for the student
type Student struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Class string `json:"class"`
	Email string `json:"email"`
	Department string `json:"department"`
	CGPA string `json:"cgpa"`
}

// Array of Student for storing student data
var students []Student

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Students API HomePage!")
	fmt.Println("Endpoint Hit: homepage")
}

func handleRequests() {
	// Handler setup
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/students", returnAllStudents)
	myRouter.HandleFunc("/student", createNewStudent).Methods("POST")
	myRouter.HandleFunc("/student/{id}", returnSingleStudent)
	myRouter.HandleFunc("/delete/{id}", deleteStudent).Methods("DELETE")
	fmt.Println("Server successfully started on port: 8080")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
	
}

func returnAllStudents(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllStudents")
	json.NewEncoder(w).Encode(students)
}

func returnSingleStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	// Loop over all of our Students
	// If the student.Id equals the key we pass in,
	// return the student encoded as JSON
	for _, student := range students {
		if student.Id == key {
			json.NewEncoder(w).Encode(student)
		}
	}
	fmt.Fprintln(w, "Student not foundwith id "+key)
}

func createNewStudent(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var stu Student
	json.Unmarshal(reqBody, &stu)
	//check for duplicates
	for _, existingStudent := range students {
		if existingStudent.Id == stu.Id {
			http.Error(w, "Student with ID "+stu.Id+" already exists", http.StatusConflict)
			return
		}
	}
	// Update our global variable of students
	// to include the new student
	students = append(students, stu)
	json.NewEncoder(w).Encode(stu)
}

func deleteStudent(w http.ResponseWriter, r *http.Request) {
	// Parsing the path parameters to read the incoming request
	vars := mux.Vars(r)

	// Extracting the 'id' of the student passed in the path
	id := vars["id"]

	for index, student := range students {
		if student.Id == id {
			students = append(students[:index], students[index+1:]...)
			fmt.Println("The student is deleted", id)
			fmt.Fprintln(w, "The Student has been deleted successfully")
			return
		}
	}

	// If student with given id is not found
	fmt.Fprintln(w, "Student not found")
}

func main() {
	// Initializing some example students
	students = []Student{
		Student{Id: "1", Name: "John Doe", Age: 20, Class: "Physics", Email:"1@gmail.com",Department:"CSE",CGPA:"8.9"},
		Student{Id: "2", Name: "Jane Smith", Age: 22, Class: "Mathematics",Email:"2@gmail.com",Department:"IT",CGPA:"9.1"},
	}

	handleRequests()
}

