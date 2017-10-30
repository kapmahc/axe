package nut

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
)

func (p *AttachmentsPlugin) currentUser(c *gin.Context) (*User, error) {
	lng := c.MustGet(web.LOCALE).(string)
	cm, err := p.Jwt.Validate([]byte(c.Query("token")))
	if err != nil {
		return nil, err
	}
	var user User
	if err = p.DB.Model(&user).Where("uid = ?", cm.Get("uid")).Limit(1).Select(); err != nil {
		return nil, err
	}
	if !user.IsConfirm() || user.IsLock() {
		return nil, p.I18n.E(lng, "errors.forbidden")
	}
	return &user, nil
}

func (p *AttachmentsPlugin) upload(c *gin.Context, name string, buf []byte, size int64) (string, error) {
	user, err := p.currentUser(c)
	if err != nil {
		return "", err
	}
	mty, url, err := p.S3.Write(name, buf, size)
	if err != nil {
		return "", err
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
		return "", err
	}
	return url, err
}

func (p *AttachmentsPlugin) list(f func(*Attachment) bool) web.UEditorManager {
	return func(c *gin.Context) ([]string, error) {
		user, err := p.currentUser(c)
		if err != nil {
			return nil, err
		}
		var items []Attachment
		err = p.DB.Model(&items).
			Column("media_type", "url").
			Where("user_id = ?", user.ID).
			Order("updated_at DESC").Select()
		if err != nil {
			return nil, err
		}

		var list []string
		if err == nil {
			for _, it := range items {
				if f(&it) {
					list = append(list, it.URL)
				}
			}
		}
		return list, nil
	}
}
