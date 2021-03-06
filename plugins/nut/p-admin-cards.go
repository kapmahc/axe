package nut

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
)

func (p *AdminPlugin) checkCardToken(user *User, tid uint) bool {
	return p.Dao.Is(user.ID, RoleAdmin)
}

func (p *AdminPlugin) editCardH(tid uint, token string) (string, string, error) {
	var it Card
	if err := p.DB.Model(&it).
		Column("id", "title", "summary").
		Where("id = ?", tid).
		Limit(1).Select(); err != nil {
		return "", "", err
	}
	return it.Title, it.Summary, nil

}
func (p *AdminPlugin) updateCardH(id uint, body string) error {
	return p.DB.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.Model(&Card{
			ID:        id,
			Summary:   body,
			Type:      web.HTML,
			UpdatedAt: time.Now(),
		}).Column("summary", "type", "updated_at").Update()
		return err
	})
}

func (p *AdminPlugin) indexCards(l string, c *gin.Context) (interface{}, error) {
	var items []Card
	err := p.DB.Model(&items).
		Where("lang = ?", l).
		Order("loc ASC").Order("sort_order ASC").Select()
	return items, err
}

func (p *AdminPlugin) showCard(l string, c *gin.Context) (interface{}, error) {
	var item Card
	err := p.DB.Model(&item).
		Where("id = ?", c.Param("id")).
		Limit(1).Select()
	return item, err
}

type fmCard struct {
	Href      string `json:"href" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Summary   string `json:"summary" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Action    string `json:"action" binding:"required"`
	Logo      string `json:"logo" binding:"required"`
	Loc       string `json:"loc" binding:"required"`
	SortOrder int    `json:"sortOrder"`
}

func (p *AdminPlugin) createCard(l string, c *gin.Context) (interface{}, error) {
	var fm fmCard
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Insert(&Card{
			Href:      fm.Href,
			Title:     fm.Title,
			Summary:   fm.Summary,
			Type:      fm.Type,
			Action:    fm.Action,
			Loc:       fm.Loc,
			Logo:      fm.Logo,
			SortOrder: fm.SortOrder,
			Lang:      l,
			UpdatedAt: time.Now(),
		})
	})
	return gin.H{}, err
}

func (p *AdminPlugin) updateCard(l string, c *gin.Context) (interface{}, error) {
	var fm fmCard
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		_, err := tx.Model(&Card{
			Href:      fm.Href,
			Title:     fm.Title,
			Summary:   fm.Summary,
			Type:      fm.Type,
			Action:    fm.Action,
			Loc:       fm.Loc,
			Logo:      fm.Logo,
			SortOrder: fm.SortOrder,
			UpdatedAt: time.Now(),
		}).
			Column("href", "title", "summary", "action", "type", "logo", "loc", "sort_order", "updated_at").
			Where("id = ?", c.Param("id")).
			Update()
		return err
	})
	return gin.H{}, err
}

func (p *AdminPlugin) destroyCard(l string, c *gin.Context) (interface{}, error) {
	_, err := p.DB.Model(&Card{}).Where("id = ?", c.Param("id")).Delete()
	return gin.H{}, err
}
