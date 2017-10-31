package nut

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
)

func (p *AttachmentsPlugin) create(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(CurrentUser).(*User)
	// c.Request.ParseMultipartForm(maxMemory)
	fd, fh, err := c.Request.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	name := fh.Filename
	size := fh.Size
	buf := make([]byte, size)
	if _, err = fd.Read(buf); err != nil {
		return nil, err
	}

	mty, url, err := p.S3.Write(name, buf, size)
	if err != nil {
		return nil, err
	}

	if err = p.DB.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Insert(&Attachment{
			Title:        name,
			Length:       size,
			MediaType:    mty,
			URL:          url,
			ResourceID:   DefaultResourceID,
			ResourceType: DefaultResourceType,
			UserID:       user.ID,
			UpdatedAt:    time.Now(),
		})
	}); err != nil {
		return nil, err
	}
	return gin.H{
		"url":    url,
		"name":   name,
		"status": "done",
		"uid":    url,
	}, err
}

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
	if admin := c.MustGet(IsAdmin).(bool); admin {
		return
	}
	user := c.MustGet(CurrentUser).(*User)
	lng := c.MustGet(web.LOCALE).(string)
	cnt, err := p.DB.Model(&Attachment{}).Where("id = ? AND user_id = ?", c.Param("id"), user.ID).Count()
	if err != nil || cnt == 0 {
		p.Layout.Abort(c, http.StatusInternalServerError, p.I18n.E(lng, "errors.forbidden"))
		return
	}

}
