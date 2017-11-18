package forum

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/plugins/nut"
	"github.com/kapmahc/axe/web"
)

func (p *Plugin) canEditComment(c *gin.Context) {
	if admin := c.MustGet(nut.IsAdmin).(bool); admin {
		return
	}
	lng := c.MustGet(web.LOCALE).(string)
	user := c.MustGet(nut.CurrentUser).(*nut.User)
	cnt, err := p.DB.Model(&Comment{}).Where("id = ? AND user_id = ?", c.Param("id"), user.ID).Count()
	if err != nil || cnt == 0 {
		p.Layout.Abort(c, http.StatusInternalServerError, p.I18n.E(lng, "errors.forbidden"))
		return
	}
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
