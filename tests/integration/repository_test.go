package integration

import (
	"book-tracker/internal/db"
	"book-tracker/internal/models"
	"book-tracker/internal/repository"
	"context"
	"testing"

	"github.com/jmoiron/sqlx"
)

func setupTestDB(t *testing.T) *sqlx.DB {
	database, err := db.NewDB()
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	return database
}

func TestCreateAndGetBook(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBookRepository(db)
	book := models.Book{Title: "Integration Test Book", Author: "Test Author", Progress: 20}

	// Test Create
	err := repo.CreateBook(context.Background(), &book)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}
	if book.ID == 0 {
		t.Error("Expected book ID to be set, got 0")
	}

	// Test Get
	books, err := repo.GetBooks(context.Background())
	if err != nil {
		t.Fatalf("Failed to get books: %v", err)
	}
	if len(books) == 0 {
		t.Error("Expected at least one book, got none")
	}
}

func TestUpdateBook(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBookRepository(db)
	book := models.Book{Title: "Test Book", Author: "Test Author", Progress: 20}

	// Create a book to update
	err := repo.CreateBook(context.Background(), &book)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	// Update the book
	book.Title = "Updated Book"
	book.Progress = 50
	err = repo.UpdateBook(context.Background(), &book)
	if err != nil {
		t.Fatalf("Failed to update book: %v", err)
	}

	// Verify update
	books, err := repo.GetBooks(context.Background())
	if err != nil {
		t.Fatalf("Failed to get books: %v", err)
	}
	found := false
	for _, b := range books {
		if b.ID == book.ID && b.Title == "Updated Book" && b.Progress == 50 {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected updated book not found")
	}
}

func TestDeleteBook(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewBookRepository(db)
	book := models.Book{Title: "Test Book", Author: "Test Author", Progress: 20}

	// Create a book to delete
	err := repo.CreateBook(context.Background(), &book)
	if err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	// Delete the book
	err = repo.DeleteBook(context.Background(), book.ID)
	if err != nil {
		t.Fatalf("Failed to delete book: %v", err)
	}

	// Verify deletion
	books, err := repo.GetBooks(context.Background())
	if err != nil {
		t.Fatalf("Failed to get books: %v", err)
	}
	for _, b := range books {
		if b.ID == book.ID {
			t.Error("Expected book to be deleted")
		}
	}
}
