package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Student struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Age        int    `json:"age"`
	Class      string `json:"class"`
	Email      string `json:"email"`
	Department string `json:"department"`
	CGPA       string `json:"cgpa"`
}


var students []Student

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Students API HomePage!")
	fmt.Println("Endpoint Hit: homepage")
}

func handleRequests() {
	// Handler setup
	myRouter := mux.NewRouter().StrictSlash(true) //mux is for directing to request to designated function

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/students", returnAllStudents)
	myRouter.HandleFunc("/student", createNewStudent).Methods("POST")
	myRouter.HandleFunc("/student/{id}", returnSingleStudent)
	myRouter.HandleFunc("/patch/{id}", patchStudent).Methods("PATCH")
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

	for _, student := range students {
		if student.Id == key {
			json.NewEncoder(w).Encode(student)
			return
		}
	}
	fmt.Fprintln(w, "Student not found with id =  "+key)
}

func createNewStudent(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var stu Student
	json.Unmarshal(reqBody, &stu) //unmarshal is for the data from reqBody is stored at stu

	//checking for duplicates
	for _, existingStudent := range students {
		if existingStudent.Id == stu.Id {
			http.Error(w, "Student with ID "+stu.Id+" already exists", http.StatusConflict)
			return
		}
	}

	students = append(students, stu)
	json.NewEncoder(w).Encode(stu)
}
func patchStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var stu *Student
	for _, student := range students {
		if student.Id == id {
			stu = &student
			break
		}
	}
	if stu == nil {
		http.Error(w, "student not found", http.StatusNotFound)
		return
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	var patchData map[string]interface{}
	if err := json.Unmarshal(reqBody, &patchData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if name, ok := patchData["name"].(string); ok {
		stu.Name = name
	}
	if age, ok := patchData["age"].(float64); ok {
		stu.Age = int(age)
	}
	if class, ok := patchData["class"].(string); ok {
		stu.Class = class
	}
	if email, ok := patchData["email"].(string); ok {
		stu.Email = email
	}
	if department, ok := patchData["department"].(string); ok {
		stu.Department = department
	}
	if cgpa, ok := patchData["cgpa"].(string); ok {
		stu.CGPA = cgpa
	}

	json.NewEncoder(w).Encode(stu)

}
func deleteStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"] //to get id

	for index, student := range students {
		if student.Id == id {
			students = append(students[:index], students[index+1:]...)
			fmt.Println("The student is deleted", id)
			fmt.Fprintln(w, "The Student has been deleted successfully")
			return
		}
	}

	fmt.Fprintln(w, "Student not found")
}

func main() {
	students = []Student{
		Student{Id: "1", Name: "John Doe", Age: 20, Class: "Physics", Email: "1@gmail.com", Department: "CSE", CGPA: "8.9"},
		Student{Id: "2", Name: "Jane Smith", Age: 22, Class: "Mathematics", Email: "2@gmail.com", Department: "IT", CGPA: "9.1"},
	}

	handleRequests()
}

