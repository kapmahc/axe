package nut

import (
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
)

func (p *AdminPlugin) showLocale(l string, c *gin.Context) (interface{}, error) {
	var item web.Locale
	err := p.DB.Model(&item).
		Where("lang = ? AND code = ?", l, c.Param("code")).
		Limit(1).Select()
	return gin.H{"code": item.Code, "message": item.Message}, err
}

type fmLocale struct {
	Code    string `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
}

func (p *AdminPlugin) createLocale(l string, c *gin.Context) (interface{}, error) {
	var fm fmLocale
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		return p.I18n.Set(tx, l, fm.Code, fm.Message)
	})
	return gin.H{}, err
}

func (p *AdminPlugin) indexLocales(l string, c *gin.Context) (interface{}, error) {
	var items []web.Locale
	err := p.DB.Model(&items).
		Column("id", "code", "message").
		Where("lang = ?", l).
		Order("code ASC").Select()
	return items, err
}

func (p *AdminPlugin) destroyLocale(l string, c *gin.Context) (interface{}, error) {
	_, err := p.DB.Model(&web.Locale{}).Where("id = ?", c.Param("id")).Delete()
	return gin.H{}, err
}
