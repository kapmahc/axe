package nut

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

func (p *AdminPlugin) indexFriendLinks(l string, c *gin.Context) (interface{}, error) {
	var items []FriendLink
	err := p.DB.Model(&items).
		Order("sort_order ASC").Select()
	return items, err
}

func (p *AdminPlugin) showFriendLink(l string, c *gin.Context) (interface{}, error) {
	var item FriendLink
	err := p.DB.Model(&item).
		Where("id = ?", c.Param("id")).
		Limit(1).Select()
	return item, err
}

type fmFriendLink struct {
	Title     string `json:"title" binding:"required"`
	Home      string `json:"home" binding:"required"`
	Logo      string `json:"logo" binding:"required"`
	SortOrder int    `json:"sortOrder"`
}

func (p *AdminPlugin) createFriendLink(l string, c *gin.Context) (interface{}, error) {
	var fm fmFriendLink
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Insert(&FriendLink{
			Home:      fm.Home,
			Title:     fm.Title,
			Logo:      fm.Logo,
			SortOrder: fm.SortOrder,
			UpdatedAt: time.Now(),
		})
	})
	return gin.H{}, err
}

func (p *AdminPlugin) updateFriendLink(l string, c *gin.Context) (interface{}, error) {
	var fm fmFriendLink
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.Model(&FriendLink{
			Home:      fm.Home,
			Title:     fm.Title,
			Logo:      fm.Logo,
			SortOrder: fm.SortOrder,
			UpdatedAt: time.Now(),
		}).
			Column("home", "title", "logo", "sort_order", "updated_at").
			Where("id = ?", c.Param("id")).
			Update()
		return err
	})
	return gin.H{}, err
}

func (p *AdminPlugin) destroyFriendLink(l string, c *gin.Context) (interface{}, error) {
	_, err := p.DB.Model(&FriendLink{}).Where("id = ?", c.Param("id")).Delete()
	return gin.H{}, err
}
