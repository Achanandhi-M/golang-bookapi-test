package unit

import (
	"book-tracker/internal/handlers"
	"book-tracker/internal/models"
	"book-tracker/internal/repository"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/gorilla/mux"
)

type mockBookRepository struct {
	createFunc func(ctx context.Context, book *models.Book) error
	getFunc    func(ctx context.Context) ([]models.Book, error)
	updateFunc func(ctx context.Context, book *models.Book) error
	deleteFunc func(ctx context.Context, id int) error
}

func (m *mockBookRepository) CreateBook(ctx context.Context, book *models.Book) error {
	return m.createFunc(ctx, book)
}

func (m *mockBookRepository) GetBooks(ctx context.Context) ([]models.Book, error) {
	return m.getFunc(ctx)
}

func (m *mockBookRepository) UpdateBook(ctx context.Context, book *models.Book) error {
	return m.updateFunc(ctx, book)
}

func (m *mockBookRepository) DeleteBook(ctx context.Context, id int) error {
	return m.deleteFunc(ctx, id)
}

var _ repository.BookRepositoryInterface = &mockBookRepository{}

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name           string
		inputBook      models.Book
		createFunc     func(ctx context.Context, b *models.Book) error
		expectedStatus int
		expectedID     int
		expectError    bool
	}{
		{
			name:           "Successful creation",
			inputBook:      models.Book{Title: "Test Book", Author: "Test Author", Progress: 10},
			createFunc:     func(ctx context.Context, b *models.Book) error { b.ID = 1; return nil },
			expectedStatus: http.StatusCreated,
			expectedID:     1,
			expectError:    false,
		},
		{
			name:           "Invalid request body",
			inputBook:      models.Book{}, // Invalid JSON
			createFunc:     func(ctx context.Context, b *models.Book) error { return nil },
			expectedStatus: http.StatusBadRequest,
			expectedID:     0,
			expectError:    true,
		},
		{
			name:           "Repository error",
			inputBook:      models.Book{Title: "Test Book", Author: "Test Author", Progress: 10},
			createFunc:     func(ctx context.Context, b *models.Book) error { return errors.New("database error") },
			expectedStatus: http.StatusInternalServerError,
			expectedID:     0,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.inputBook)
			req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			mockRepo := &mockBookRepository{
				createFunc: tt.createFunc,
			}
			router := mux.NewRouter()
			router.HandleFunc("/books", handlers.CreateBook(mockRepo)).Methods("POST")
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				var createdBook models.Book
				if err := json.NewDecoder(w.Body).Decode(&createdBook); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if createdBook.ID != tt.expectedID {
					t.Errorf("Expected book ID %d, got %d", tt.expectedID, createdBook.ID)
				}
				if createdBook.Title != tt.inputBook.Title {
					t.Errorf("Expected book title %s, got %s", tt.inputBook.Title, createdBook.Title)
				}
			}
		})
	}
}

func TestGetBooks(t *testing.T) {
	tests := []struct {
		name           string
		getFunc        func(ctx context.Context) ([]models.Book, error)
		expectedStatus int
		expectedBooks  []models.Book
	}{
		{
			name: "Successful get",
			getFunc: func(ctx context.Context) ([]models.Book, error) {
				return []models.Book{{ID: 1, Title: "Test Book", Author: "Test Author"}}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBooks:  []models.Book{{ID: 1, Title: "Test Book", Author: "Test Author"}},
		},
		{
			name:           "Repository error",
			getFunc:        func(ctx context.Context) ([]models.Book, error) { return nil, errors.New("database error") },
			expectedStatus: http.StatusInternalServerError,
			expectedBooks:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/books", nil)
			w := httptest.NewRecorder()

			mockRepo := &mockBookRepository{
				getFunc: tt.getFunc,
			}
			router := mux.NewRouter()
			router.HandleFunc("/books", handlers.GetBooks(mockRepo)).Methods("GET")
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBooks != nil {
				var books []models.Book
				if err := json.NewDecoder(w.Body).Decode(&books); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if len(books) != len(tt.expectedBooks) {
					t.Errorf("Expected %d books, got %d", len(tt.expectedBooks), len(books))
				}
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		inputBook      models.Book
		updateFunc     func(ctx context.Context, book *models.Book) error
		expectedStatus int
		expectError    bool
	}{
		{
			name:           "Successful update",
			id:             "1",
			inputBook:      models.Book{Title: "Updated Book", Author: "Updated Author", Progress: 20},
			updateFunc:     func(ctx context.Context, book *models.Book) error { return nil },
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid ID",
			id:             "invalid",
			inputBook:      models.Book{Title: "Updated Book", Author: "Updated Author"},
			updateFunc:     func(ctx context.Context, book *models.Book) error { return nil },
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Repository error",
			id:             "1",
			inputBook:      models.Book{Title: "Updated Book", Author: "Updated Author"},
			updateFunc:     func(ctx context.Context, book *models.Book) error { return errors.New("database error") },
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.inputBook)
			req := httptest.NewRequest(http.MethodPut, "/books/"+tt.id, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			mockRepo := &mockBookRepository{
				updateFunc: tt.updateFunc,
			}
			router := mux.NewRouter()
			router.HandleFunc("/books/{id}", handlers.UpdateBook(mockRepo)).Methods("PUT")
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				var updatedBook models.Book
				if err := json.NewDecoder(w.Body).Decode(&updatedBook); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if updatedBook.Title != tt.inputBook.Title {
					t.Errorf("Expected book title %s, got %s", tt.inputBook.Title, updatedBook.Title)
				}
			}
		})
	}
}

func TestDeleteBook(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		deleteFunc     func(ctx context.Context, id int) error
		expectedStatus int
	}{
		{
			name:           "Successful delete",
			id:             "1",
			deleteFunc:     func(ctx context.Context, id int) error { return nil },
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Invalid ID",
			id:             "invalid",
			deleteFunc:     func(ctx context.Context, id int) error { return nil },
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Repository error",
			id:             "1",
			deleteFunc:     func(ctx context.Context, id int) error { return errors.New("database error") },
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/books/"+tt.id, nil)
			w := httptest.NewRecorder()

			mockRepo := &mockBookRepository{
				deleteFunc: tt.deleteFunc,
			}
			router := mux.NewRouter()
			router.HandleFunc("/books/{id}", handlers.DeleteBook(mockRepo)).Methods("DELETE")
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
