package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//UserDefinition a custom word and body of text defining what the word means
type UserDefinition struct {
	Word   string  `json:"word"`
	Body   string  `json:"body"`
	Author *Author `json:"author"`
}

//Author the creator of an User Definition
type Author struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

//mock collection
var userDefinitions []UserDefinition

func createDefinition(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application.json")
	var userDefinition UserDefinition
	_ = json.NewDecoder(request.Body).Decode(&userDefinition)
	userDefinitions = append(userDefinitions, userDefinition)
	json.NewEncoder(response).Encode(userDefinition)
}
func getAllWords(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application.json")
	json.NewEncoder(response).Encode(userDefinitions)
}
func updateDefinition(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application.json")
	params := mux.Vars(request)
	log.Println(params)
	for index, userDef := range userDefinitions {
		if userDef.Word == params["word"] {

			log.Println("found word : " + userDef.Word)
			var userDefinition UserDefinition
			_ = json.NewDecoder(request.Body).Decode(&userDefinition)

			userDefinitions[index].Body = userDefinition.Body
			userDefinitions[index].Author = userDefinition.Author
			return
		}
	}
}
func getSingleDefinition(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application.json")
	params := mux.Vars(request)
	for _, userDef := range userDefinitions {
		if userDef.Word == params["word"] {
			json.NewEncoder(response).Encode(userDef)
		}
	}
}
func removeDefinition(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application.json")
	params := mux.Vars(request)
	for index, userDef := range userDefinitions {
		if userDef.Word == params["word"] {
			userDefinitions = append(userDefinitions[:index], userDefinitions[index+1:]...)
			break
		}
	}
}
func getAllAuthors(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application.json")
	var authors = make(map[string]string)

	for _, userDef := range userDefinitions {
		authors[userDef.Author.Id] = userDef.Author.Name
	}
	json.NewEncoder(response).Encode(authors)

}
func getAllWordsByAuthor(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application.json")
	params := mux.Vars(request)
	var words []string

	for _, userDef := range userDefinitions {
		if userDef.Author.Id == params["author"] {
			words = append(words, userDef.Word)
		}

	}
	json.NewEncoder(response).Encode(words)
}

func main() {
	//Initializing mux router
	router := mux.NewRouter()

	//MOCK Data
	userDefinitions = append(userDefinitions, UserDefinition{Word: "breakfast", Body: "The most important meal of the day", Author: &Author{Id: "1", Name: "chris"}})
	userDefinitions = append(userDefinitions, UserDefinition{Word: "dinner", Body: "The most important meal of the day", Author: &Author{Id: "1", Name: "chris"}})

	//Route handlers / paths
	router.HandleFunc("/api/definitions", createDefinition).Methods("POST")
	router.HandleFunc("/api/definitions", getAllWords).Methods("GET")
	router.HandleFunc("/api/definitions/{word}", updateDefinition).Methods("PUT")
	router.HandleFunc("/api/definitions/{word}", getSingleDefinition).Methods("GET")
	router.HandleFunc("/api/definitions/{word}", removeDefinition).Methods("DELETE")
	router.HandleFunc("/api/authors", getAllAuthors).Methods("GET")
	router.HandleFunc("/api/authors/{author}", getAllWordsByAuthor).Methods("GET")

	log.Println(http.ListenAndServe(":8083", router))
}
