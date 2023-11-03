package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book Struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Init books var as a slice Book struct
var books []Book

// Get All Books
func getBooks(w http.ResponseWriter, r *http.Request) { // w is writer, r is request
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		log.Fatalln("There was an error encoding the initialized struct")
	}
}

// Get Single Book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// Loop through books and find book with ID
	for _, book := range books {
		if book.ID == params["id"] {
			err := json.NewEncoder(w).Encode(book)
			if err != nil {
				log.Fatalln("There was an error encoding the initialized struct")
			}
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Create a New Book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}
	book.ID = strconv.Itoa(rand.Intn(10000000)) // ID just for example, this can generate the same id twice, so not to be used in production
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// Update Book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// Loop through books and find book with ID
	for index, book := range books {
		if book.ID == params["id"] {
			// delete book, filter out book
			books = append(books[:index], books[index+1:]...)
			var updatedBook Book
			err := json.NewDecoder(r.Body).Decode(&updatedBook)
			if err != nil {
				log.Fatalln("There was an error decoding the request body into the struct")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			updatedBook.ID = params["id"]
			books = append(books, updatedBook)
			json.NewEncoder(w).Encode(updatedBook)
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

// Delete Book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// Loop through books and find book with ID
	for index, book := range books {
		if book.ID == params["id"] {
			// delete book, filter out book
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

// Index route "/"
func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}

func main() {

	// Init Router
	r := mux.NewRouter().UseEncodedPath()

	// Mock Data
	books = append(books, Book{ID: "1", Isbn: "12345", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Taylor"}})
	books = append(books, Book{ID: "2", Isbn: "776655", Title: "Book Two", Author: &Author{Firstname: "Tim", Lastname: "Banks"}})
	books = append(books, Book{ID: "3", Isbn: "45243", Title: "Book Three", Author: &Author{Firstname: "Ray", Lastname: "Smith"}})

	// Route Handlers / Endpoints
	r.HandleFunc("/", index)
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")
	fmt.Println("Server starting...")
	log.Fatal(http.ListenAndServe(":8000", r))
}
