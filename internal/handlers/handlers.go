package handlers

import (
	"book-tracker/internal/models"
	"book-tracker/internal/repository"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func RegisterBookHandlers(router *mux.Router, db *sqlx.DB) {
	repo := repository.NewBookRepository(db)
	router.HandleFunc("/books", CreateBook(repo)).Methods("POST")
	router.HandleFunc("/books", GetBooks(repo)).Methods("GET")
	router.HandleFunc("/books/{id}", UpdateBook(repo)).Methods("PUT")
	router.HandleFunc("/books/{id}", DeleteBook(repo)).Methods("DELETE")
}

func CreateBook(repo repository.BookRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var book models.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		// Validate book input
		if book.Title == "" || book.Author == "" || book.Progress < 0 {
			http.Error(w, "Title, Author, and non-negative Progress are required", http.StatusBadRequest)
			return
		}
		if err := repo.CreateBook(r.Context(), &book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(book)
	}
}

func GetBooks(repo repository.BookRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := repo.GetBooks(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(books)
	}
}

func UpdateBook(repo repository.BookRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		var book models.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		// Validate book input
		if book.Title == "" || book.Author == "" || book.Progress < 0 {
			http.Error(w, "Title, Author, and non-negative Progress are required", http.StatusBadRequest)
			return
		}
		book.ID = id
		if err := repo.UpdateBook(r.Context(), &book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(book)
	}
}

func DeleteBook(repo repository.BookRepositoryInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		if err := repo.DeleteBook(r.Context(), id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
