package forum

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/orm"
	"github.com/kapmahc/axe/plugins/nut"
)

func (p *Plugin) indexArticlesH(l string, d gin.H, c *gin.Context) error {
	var items []Article
	if err := p.DB.Model(&items).
		Column("id", "title", "body", "type", "updated_at").
		Order("updated_at DESC").
		Select(); err != nil {
		return err
	}
	d["items"] = items
	d[nut.TITLE] = p.I18n.T(l, "forum.articles.index.title")
	return nil
}

func (p *Plugin) indexTagsH(l string, d gin.H, c *gin.Context) error {
	var items []Tag
	if err := p.DB.Model(&items).
		Column("id", "name").
		Order("updated_at DESC").
		Select(); err != nil {
		return err
	}
	d["items"] = items
	d[nut.TITLE] = p.I18n.T(l, "forum.tags.index.title")
	return nil
}

func (p *Plugin) indexCommentsH(l string, d gin.H, c *gin.Context) error {
	var items []Comment
	if err := p.DB.Model(&items).
		Column("id", "article_id", "body", "type", "updated_at").
		Order("updated_at DESC").
		Select(); err != nil {
		return err
	}
	d["items"] = items
	d[nut.TITLE] = p.I18n.T(l, "forum.comments.index.title")
	return nil
}

func (p *Plugin) showArticleH(l string, d gin.H, c *gin.Context) error {
	var item Article
	if err := p.DB.Model(&item).
		Where("id = ?", c.Param("id")).
		Relation("Tags", func(q *orm.Query) (*orm.Query, error) {
			return q.Column("id", "name").Order("updated_at DESC"), nil
		}).Relation("Comments", func(q *orm.Query) (*orm.Query, error) {
		return q.Column("body", "type").Order("updated_at ASC"), nil
	}).
		Limit(1).Select(); err != nil {
		return err
	}
	d["item"] = item
	d[nut.TITLE] = item.Title
	return nil
}

func (p *Plugin) showTagH(l string, d gin.H, c *gin.Context) error {
	var item Tag
	if err := p.DB.Model(&item).
		Where("id = ?", c.Param("id")).
		Relation("Articles", func(q *orm.Query) (*orm.Query, error) {
			return q.Column("id", "title", "body", "type", "updated_at"), nil
		}).
		Limit(1).Select(); err != nil {
		return err
	}
	d["item"] = item
	d[nut.TITLE] = item.Name
	return nil
}
