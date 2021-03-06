package forum

import (
	"time"

	"github.com/kapmahc/axe/plugins/nut"
)

// Article article
type Article struct {
	tableName struct{}  `sql:"forum_articles"`
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	User      nut.User  `json:"user"`
	UserID    uint      `json:"userId"`
	Tags      []Tag     `json:"tags" pg:",many2many:forum_articles_tags"`
	Comments  []Comment `json:"comments"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// ArticleTag articles-tags
type ArticleTag struct {
	tableName struct{} `sql:"forum_articles_tags"`
	ArticleID uint
	TagID     uint
}

// Tag tag
type Tag struct {
	tableName struct{}  `sql:"forum_tags"`
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Articles  []Article `json:"articles" pg:",many2many:forum_articles_tags"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// Comment comment
type Comment struct {
	tableName struct{}  `sql:"forum_comments"`
	ID        uint      `json:"id"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	User      nut.User  `json:"user"`
	UserID    uint      `json:"userId"`
	Article   Article   `json:"article"`
	ArticleID uint      `json:"articleId"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}
