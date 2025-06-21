package api

import (
	"book-tracker/internal/db"
	"book-tracker/internal/handlers"
	"book-tracker/internal/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

func setupTestServer(t *testing.T) *mux.Router {
	db, err := db.NewDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	router := mux.NewRouter()
	handlers.RegisterBookHandlers(router, db)
	return router
}

func TestAPICreateAndGet(t *testing.T) {
	router := setupTestServer(t)

	// Test Create
	book := models.Book{Title: "API Test Book", Author: "API Author", Progress: 30}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var createdBook models.Book
	if err := json.NewDecoder(w.Body).Decode(&createdBook); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Test Get
	req = httptest.NewRequest(http.MethodGet, "/books", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var books []models.Book
	if err := json.NewDecoder(w.Body).Decode(&books); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(books) == 0 {
		t.Error("Expected at least one book, got none")
	}
}

func TestAPIUpdateBook(t *testing.T) {
	router := setupTestServer(t)

	// Create a book to update
	book := models.Book{Title: "API Test Book", Author: "API Author", Progress: 30}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var createdBook models.Book
	if err := json.NewDecoder(w.Body).Decode(&createdBook); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Update the book
	updatedBook := models.Book{Title: "Updated API Book", Author: "Updated Author", Progress: 50}
	body, _ = json.Marshal(updatedBook) // Fixed: Use = instead of := to reuse body
	req = httptest.NewRequest(http.MethodPut, "/books/"+strconv.Itoa(createdBook.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resultBook models.Book
	if err := json.NewDecoder(w.Body).Decode(&resultBook); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resultBook.Title != updatedBook.Title {
		t.Errorf("Expected book title %s, got %s", updatedBook.Title, resultBook.Title)
	}
}

func TestAPIDeleteBook(t *testing.T) {
	router := setupTestServer(t)

	// Create a book to delete
	book := models.Book{Title: "API Test Book", Author: "API Author", Progress: 30}
	body, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var createdBook models.Book
	if err := json.NewDecoder(w.Body).Decode(&createdBook); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Delete the book
	req = httptest.NewRequest(http.MethodDelete, "/books/"+strconv.Itoa(createdBook.ID), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify deletion
	req = httptest.NewRequest(http.MethodGet, "/books", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var books []models.Book
	if err := json.NewDecoder(w.Body).Decode(&books); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	for _, b := range books {
		if b.ID == createdBook.ID {
			t.Error("Expected book to be deleted")
		}
	}
}
