package forum

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/plugins/nut"
	"github.com/kapmahc/axe/web"
)

func (p *Plugin) checkCommentToken(user *nut.User, cid uint) bool {
	var it Comment
	if err := p.DB.Model(&it).Column("user_id").Where("id = ?", cid).Limit(1).Select(); err != nil {
		return false
	}
	return it.UserID == user.ID || p.Dao.Is(user.ID, nut.RoleAdmin)
}

func (p *Plugin) editCommentH(tid uint, token string) (string, string, string, error) {
	var it Comment
	if err := p.DB.Model(&it).
		Column("id", "body").
		Where("id = ?", tid).
		Limit(1).Select(); err != nil {
		return "", "", "", err
	}
	return strconv.Itoa(int(it.ID)), fmt.Sprintf("/forum/comments/edit/%s", token), it.Body, nil

}
func (p *Plugin) updateCommentH(id uint, body string) error {
	return p.DB.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.Model(&Comment{
			ID:        id,
			Body:      body,
			Type:      web.HTML,
			UpdatedAt: time.Now(),
		}).Column("body", "type", "updated_at").Update()
		return err
	})
}

func (p *Plugin) indexComments(l string, c *gin.Context) (interface{}, error) {
	var items []Comment
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	admin := c.MustGet(nut.IsAdmin).(bool)
	db := p.DB.Model(&items).Column("id", "body", "type")
	if !admin {
		db = db.Where("user_id = ?", user.ID)
	}
	err := db.Order("updated_at DESC").Select()
	return items, err
}

func (p *Plugin) showComment(l string, c *gin.Context) (interface{}, error) {
	var item Comment
	err := p.DB.Model(&item).
		Where("id = ?", c.Param("id")).
		Limit(1).Select()
	return item, err
}

type fmComment struct {
	Body string `json:"body" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (p *Plugin) createComment(l string, c *gin.Context) (interface{}, error) {
	aid, err := strconv.Atoi(c.Query("articleId"))
	if err != nil {
		return nil, err
	}
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	var fm fmComment
	if err = c.BindJSON(&fm); err != nil {
		return nil, err
	}
	err = p.DB.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Insert(&Comment{
			Body:      fm.Body,
			Type:      fm.Type,
			ArticleID: uint(aid),
			UserID:    user.ID,
			UpdatedAt: time.Now(),
		})
	})
	return gin.H{}, err
}

func (p *Plugin) updateComment(l string, c *gin.Context) (interface{}, error) {
	var fm fmComment
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.Model(&Comment{
			Type:      fm.Type,
			Body:      fm.Body,
			UpdatedAt: time.Now(),
		}).
			Column("body", "type", "updated_at").
			Where("id = ?", c.Param("id")).
			Update()
		return err
	})
	return gin.H{}, err
}

func (p *Plugin) destroyComment(l string, c *gin.Context) (interface{}, error) {
	_, err := p.DB.Model(&Comment{}).Where("id = ?", c.Param("id")).Delete()
	return gin.H{}, err
}
