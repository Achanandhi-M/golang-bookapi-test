package models

type Book struct {
	ID       int    `json:"id" db:"id"`
	Title    string `json:"title" db:"title"`
	Author   string `json:"author" db:"author"`
	Progress int    `json:"progress" db:"progress"`
	Notes    string `json:"notes" db:"notes"`
	Finished bool   `json:"finished" db:"finished"`
	Rating   int    `json:"rating" db:"rating"`
}
