package reading

import (
	"time"

	"github.com/kapmahc/axe/plugins/nut"
)

// Book book
type Book struct {
	tableName   struct{} `sql:"reading_books"`
	ID          uint
	Author      string
	Publisher   string
	Title       string
	Type        string
	Lang        string
	File        string
	Subject     string
	Description string
	PublishedAt time.Time
	Cover       string
	UpdatedAt   time.Time
	CreatedAt   time.Time
}

// Note note
type Note struct {
	tableName struct{} `sql:"reading_notes"`
	ID        uint
	Type      string
	Body      string
	User      nut.User
	Book      Book
}
