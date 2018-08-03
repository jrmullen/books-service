package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Initialize books slice
var books []Book

func main() {
	// Initialize router
	r := mux.NewRouter()

	// Hard-coded book data
	books = append(books, Book{
		"1",
		"978-1-891830-75-4",
		"First Title",
		&Author{
			"Joe",
			"Schmo",
		}})

	books = append(books, Book{
		"2",
		"978-1-60309-050-6",
		"Second Title",
		&Author{
			"Jane",
			"Smith",
		}})

	// Endpoints
	r.HandleFunc("/api/v1/books", getBooks).Methods("GET")
	r.HandleFunc("/api/v1/book/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/v1/book", createBook).Methods("POST")
	r.HandleFunc("/api/v1/book/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/v1/book/{id}", deleteBook).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8000", r))
}

// Delete a single book by id. Returns the resulting list of books
func deleteBook(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request) // Get the parameters of the request

	for index, book := range books {
		if book.Id == params["id"] {
			books = append(books[:index], books[index+1:]...)
			json.NewEncoder(writer).Encode(books)
			return
		}
	}

	// 404 if id does not exist
	http.NotFound(writer, request)
}

// Update a book by id
func updateBook(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)

	for index, book := range books {
		if book.Id == params["id"] {
			// Delete the old book
			books = append(books[:index], books[index+1:]...)

			// Create a new book from the request body json
			var book Book
			_ = json.NewDecoder(request.Body).Decode(&book)
			book.Id = params["id"]

			// Add the newly updated book to the slice of books
			books = append(books, book)

			// Return the new book object in the body
			json.NewEncoder(writer).Encode(book)

			return
		}
	}

	// 404 if id does not exist
	http.NotFound(writer, request)
}

func createBook(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var book Book

	// Encode the request body JSON to book struct
	_ = json.NewDecoder(request.Body).Decode(&book)

	// Append the new book to the books slice
	books = append(books, book)

	// Return the new book struct
	json.NewEncoder(writer).Encode(book)
}

// Return a single book by id
func getBook(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request) // Get the parameters of the request

	// Loop through the list of books and return the element with matching id
	for _, book := range books {
		if book.Id == params["id"] {
			json.NewEncoder(writer).Encode(book)
			return
		}
	}

	// 404 if id does not exist
	http.NotFound(writer, request)
}

// Return all books
func getBooks(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(books)
}
