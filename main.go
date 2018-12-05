package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

// Book Struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author Struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Init books var as a slice Book struct
var books []Book

// Functions like a ternary operator for two strings
func ternaryForString(a, b string) (chosen string) {
	if a == "" {
		return b
	}
	return a
}

// Creates a new book with new data, uses old book data if none present
func mergeBooks(oldBook Book, newBook Book) (mergedBook Book) {
	newBook.ID = oldBook.ID
	newBook.Isbn = ternaryForString(newBook.Isbn, oldBook.Isbn)
	newBook.Title = ternaryForString(newBook.Title, oldBook.Title)
	newBook.Author.Firstname = ternaryForString(newBook.Author.Firstname, oldBook.Author.Firstname)
	newBook.Author.Lastname = ternaryForString(newBook.Author.Lastname, oldBook.Author.Lastname)
	return newBook
}

// Get All Books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get Single Book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Get params
	// Loop through books and find correct Id
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Create a New Book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(1000000)) // Mock ID not safe
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// Update a Book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			oldBook := books[index]
			books = append(books[:index], books[index+1:]...)

			var newBook Book
			_ = json.NewDecoder(r.Body).Decode(&newBook)
			newBook = mergeBooks(oldBook, newBook)

			books = append(books, newBook)
			json.NewEncoder(w).Encode(newBook)
			return
		}
	}
	json.NewEncoder(w).Encode(books)
}

// Delete a Book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	// Init Router
	r := mux.NewRouter()

	// Mock Data - @todo - implement DB
	books = append(books,
		Book{
			ID:     "1",
			Isbn:   "448743",
			Title:  "Book One",
			Author: &Author{Firstname: "John", Lastname: "Doe"},
		})

	books = append(books,
		Book{
			ID:     "2",
			Isbn:   "847564",
			Title:  "Book Two",
			Author: &Author{Firstname: "Steve", Lastname: "Smith"},
		})

	// Route Handlers / Endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", r))
}
