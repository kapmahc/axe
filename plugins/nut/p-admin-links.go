package nut

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

func (p *AdminPlugin) indexLinks(l string, c *gin.Context) (interface{}, error) {
	var items []Link
	err := p.DB.Model(&items).
		Where("lang = ?", l).
		Order("loc ASC").Order("sort_order ASC").Select()
	return items, err
}

func (p *AdminPlugin) showLink(l string, c *gin.Context) (interface{}, error) {
	var item Link
	err := p.DB.Model(&item).
		Where("id = ?", c.Param("id")).
		Limit(1).Select()
	return item, err
}

type fmLink struct {
	Href      string `json:"href" binding:"required"`
	Label     string `json:"label" binding:"required"`
	Loc       string `json:"loc" binding:"required"`
	SortOrder int    `json:"sortOrder"`
}

func (p *AdminPlugin) createLink(l string, c *gin.Context) (interface{}, error) {
	var fm fmLink
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Insert(&Link{
			Href:      fm.Href,
			Label:     fm.Label,
			Loc:       fm.Loc,
			SortOrder: fm.SortOrder,
			Lang:      l,
			UpdatedAt: time.Now(),
		})
	})
	return gin.H{}, err
}

func (p *AdminPlugin) updateLink(l string, c *gin.Context) (interface{}, error) {
	var fm fmLink
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.Model(&Link{
			Href:      fm.Href,
			Label:     fm.Label,
			Loc:       fm.Loc,
			SortOrder: fm.SortOrder,
			Lang:      l,
			UpdatedAt: time.Now(),
		}).
			Column("href", "label", "loc", "sort_order", "lang", "updated_at").
			Where("id = ?", c.Param("id")).
			Update()
		return err
	})
	return gin.H{}, err
}

func (p *AdminPlugin) destroyLink(l string, c *gin.Context) (interface{}, error) {
	_, err := p.DB.Model(&Link{}).Where("id = ?", c.Param("id")).Delete()
	return gin.H{}, err
}
