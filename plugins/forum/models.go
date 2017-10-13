package forum

import (
	"time"

	"github.com/kapmahc/axe/plugins/nut"
)

// Article article
type Article struct {
	tableName struct{} `sql:"forum_articles"`
	ID        uint
	Title     string
	Body      string
	Type      string
	User      nut.User
	Tags      []Tag
	Comments  []Comment
	UpdatedAt time.Time
	CreatedAt time.Time
}

// Tag tag
type Tag struct {
	tableName struct{} `sql:"forum_tags"`
	ID        uint
	Name      string
	Articles  []Article
	UpdatedAt time.Time
	CreatedAt time.Time
}

// Comment comment
type Comment struct {
	tableName struct{} `sql:"forum_comments"`
	ID        uint
	Body      string
	Type      string
	User      nut.User
	Article   Article
	UpdatedAt time.Time
	CreatedAt time.Time
}
