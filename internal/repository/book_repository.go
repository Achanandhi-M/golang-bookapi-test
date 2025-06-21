package repository

import (
	"book-tracker/internal/models"
	"context"

	"github.com/jmoiron/sqlx"
)

// BookRepositoryInterface defines the methods for book repository operations.
type BookRepositoryInterface interface {
	CreateBook(ctx context.Context, book *models.Book) error
	GetBooks(ctx context.Context) ([]models.Book, error)
	UpdateBook(ctx context.Context, book *models.Book) error
	DeleteBook(ctx context.Context, id int) error
}

type BookRepository struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) *BookRepository {
	return &BookRepository{db: db}
}

// Ensure BookRepository implements BookRepositoryInterface
var _ BookRepositoryInterface = &BookRepository{}

func (r *BookRepository) CreateBook(ctx context.Context, book *models.Book) error {
	query := `
		INSERT INTO books (title, author, progress, notes, finished, rating)
		VALUES (:title, :author, :progress, :notes, :finished, :rating)
		RETURNING id`
	rows, err := r.db.NamedQueryContext(ctx, query, book)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&book.ID)
	}
	return nil
}

func (r *BookRepository) GetBooks(ctx context.Context) ([]models.Book, error) {
	var books []models.Book
	query := `SELECT id, title, author, progress, notes, finished, rating FROM books`
	err := r.db.SelectContext(ctx, &books, query)
	return books, err
}

func (r *BookRepository) UpdateBook(ctx context.Context, book *models.Book) error {
	query := `
		UPDATE books
		SET title = :title, author = :author, progress = :progress, notes = :notes,
		    finished = :finished, rating = :rating
		WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, book)
	return err
}

func (r *BookRepository) DeleteBook(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
