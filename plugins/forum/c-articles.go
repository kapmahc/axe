package forum

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/kapmahc/axe/plugins/nut"
	"github.com/kapmahc/axe/web"
)

func (p *Plugin) checkArticleToken(user *nut.User, aid uint) bool {
	var it Article
	if err := p.DB.Model(&it).Column("user_id").Where("id = ?", aid).Limit(1).Select(); err != nil {
		return false
	}
	return it.UserID == user.ID || p.Dao.Is(user.ID, nut.RoleAdmin)
}

func (p *Plugin) editArticleH(tid uint, token string) (string, string, string, error) {
	var it Article
	if err := p.DB.Model(&it).
		Column("id", "title", "body").
		Where("id = ?", tid).
		Limit(1).Select(); err != nil {
		return "", "", "", err
	}
	return it.Title, fmt.Sprintf("/forum/articles/edit/%s", token), it.Body, nil

}
func (p *Plugin) updateArticleH(id uint, body string) error {
	return p.DB.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.Model(&Article{
			ID:        id,
			Body:      body,
			Type:      web.HTML,
			UpdatedAt: time.Now(),
		}).Column("body", "type", "updated_at").Update()
		return err
	})
}

func (p *Plugin) indexArticles(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	admin := c.MustGet(nut.IsAdmin).(bool)

	var items []Article
	db := p.DB.Model(&items).Column("id", "title")
	if !admin {
		db = db.Where("user_id = ?", user.ID)
	}
	err := db.Order("updated_at DESC").Select()
	return items, err
}

func (p *Plugin) showArticle(l string, c *gin.Context) (interface{}, error) {
	db := p.DB
	var item Article
	err := db.Model(&item).
		Where("id = ?", c.Param("id")).Relation("Tags", func(q *orm.Query) (*orm.Query, error) {
		return q.Column("id", "name"), nil
	}).
		Limit(1).Select()
	if err != nil {
		return nil, err
	}

	return item, nil
}

type fmArticle struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
	Type  string `json:"type" binding:"required"`
	Tags  []uint `json:"tags"`
}

func (p *Plugin) createArticle(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var fm fmArticle
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}

	err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		it := Article{
			Title:     fm.Title,
			Body:      fm.Body,
			Type:      fm.Type,
			UserID:    user.ID,
			UpdatedAt: time.Now(),
		}
		if err := tx.Insert(&it); err != nil {
			return err
		}
		var ats []interface{}
		for _, t := range fm.Tags {
			ats = append(ats, &ArticleTag{TagID: t, ArticleID: it.ID})
		}
		if err := tx.Insert(ats...); err != nil {
			return err
		}

		return nil
	})
	return gin.H{}, err
}

func (p *Plugin) updateArticle(l string, c *gin.Context) (interface{}, error) {
	var fm fmArticle
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	aid, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, err
	}
	err = p.DB.RunInTransaction(func(tx *pg.Tx) error {
		_, er := tx.Model(&Article{
			Title:     fm.Title,
			Body:      fm.Body,
			Type:      fm.Type,
			UpdatedAt: time.Now(),
		}).
			Column("title", "body", "type", "updated_at").
			Where("id = ?", aid).
			Update()
		if er != nil {
			return er
		}
		if _, er := tx.Model(&ArticleTag{}).Where("article_id = ?", aid).Delete(); er != nil {
			return er
		}
		var ats []interface{}
		for _, t := range fm.Tags {
			ats = append(ats, &ArticleTag{TagID: t, ArticleID: uint(aid)})
		}
		if er := tx.Insert(ats...); er != nil {
			return er
		}
		return nil
	})
	return gin.H{}, err
}

func (p *Plugin) destroyArticle(l string, c *gin.Context) (interface{}, error) {
	_, err := p.DB.Model(&Article{}).Where("id = ?", c.Param("id")).Delete()
	return gin.H{}, err
}
