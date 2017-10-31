package nut

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kapmahc/axe/web"
)

func (p *AttachmentsPlugin) index(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(CurrentUser).(*User)
	admin := c.MustGet(IsAdmin).(bool)
	var items []Attachment
	db := p.DB.Model(&items).Column("id", "title", "url", "media_type", "length", "updated_at")
	if !admin {
		db = db.Where("user_id = ?", user.ID)
	}
	err := db.Order("updated_at DESC").Select()
	return items, err
}

func (p *AttachmentsPlugin) destroy(l string, c *gin.Context) (interface{}, error) {
	_, err := p.DB.Model(&Attachment{}).Where("id = ?", c.Param("id")).Delete()
	return gin.H{}, err
}

func (p *AttachmentsPlugin) canEdit(c *gin.Context) {
	user := c.MustGet(CurrentUser).(*User)
	admin := c.MustGet(IsAdmin).(bool)
	lng := c.MustGet(web.LOCALE).(string)
	if !admin {
		cnt, err := p.DB.Model(&Attachment{}).Where("id = ? AND user_id = ?", c.Param("id"), user.ID).Count()
		if err != nil || cnt == 0 {
			p.Layout.Abort(c, http.StatusInternalServerError, p.I18n.E(lng, "errors.forbidden"))
			return
		}
	}
}
