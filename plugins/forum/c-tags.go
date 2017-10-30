package forum

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

func (p *Plugin) indexTags(l string, c *gin.Context) (interface{}, error) {
	var items []Tag
	err := p.DB.Model(&items).Column("id", "name").
		Order("updated_at DESC").Select()
	return items, err
}

func (p *Plugin) showTag(l string, c *gin.Context) (interface{}, error) {
	var item Tag
	err := p.DB.Model(&item).
		Where("id = ?", c.Param("id")).
		Limit(1).Select()
	return item, err
}

type fmTag struct {
	Name string `json:"name" binding:"required"`
}

func (p *Plugin) createTag(l string, c *gin.Context) (interface{}, error) {
	var fm fmTag
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Insert(&Tag{
			Name:      fm.Name,
			UpdatedAt: time.Now(),
		})
	})
	return gin.H{}, err
}

func (p *Plugin) updateTag(l string, c *gin.Context) (interface{}, error) {
	var fm fmTag
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.Model(&Tag{
			Name:      fm.Name,
			UpdatedAt: time.Now(),
		}).
			Column("name", "updated_at").
			Where("id = ?", c.Param("id")).
			Update()
		return err
	})
	return gin.H{}, err
}

func (p *Plugin) destroyTag(l string, c *gin.Context) (interface{}, error) {
	_, err := p.DB.Model(&Tag{}).Where("id = ?", c.Param("id")).Delete()
	return gin.H{}, err
}
