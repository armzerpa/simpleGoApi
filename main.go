package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type Person struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

var (
	people []Person
	mutex  sync.Mutex
	nextID int
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/person", createPerson).Methods("POST")
	router.HandleFunc("/person/{id}", getPerson).Methods("GET")
	router.HandleFunc("/people", getAllPeople).Methods("GET")

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	var person Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	person.ID = nextID
	nextID++
	people = append(people, person)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	for _, person := range people {
		if person.ID == id {
			json.NewEncoder(w).Encode(person)
			return
		}
	}

	http.Error(w, "Person not found", http.StatusNotFound)
}

func getAllPeople(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	json.NewEncoder(w).Encode(people)
}
