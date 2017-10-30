package nut

import (
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

func (p *HomePlugin) getHome(l string, d gin.H, c *gin.Context) error {
	// TODO
	return nil
}

type fmToken struct {
	Act string `json:"act" binding:"required"`
	Tid uint   `json:"tid" binding:"required"`
}

func (p *HomePlugin) postToken(l string, c *gin.Context) (interface{}, error) {
	user := c.MustGet(CurrentUser).(*User)
	var fm fmToken
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	token, err := p.Jwt.Sum(jws.Claims{
		"act": fm.Act,
		"tid": fm.Tid,
		"uid": user.UID,
	}, time.Hour*3)
	if err != nil {
		return nil, err
	}
	return gin.H{
		"token": string(token),
	}, nil
}

func (p *HomePlugin) getSiteInfo(l string, c *gin.Context) (interface{}, error) {
	// -----------
	langs, err := p.I18n.Languages()
	if err != nil {
		return nil, err
	}
	data := gin.H{"locale": l, "languages": langs}
	// -----------
	for _, k := range []string{"title", "subhead", "keywords", "description", "copyright"} {
		data[k] = p.I18n.T(l, "site."+k)
	}
	// -----------
	var author map[string]string
	if err := p.Settings.Get("site.author", &author); err != nil {
		author = map[string]string{
			"name":  "",
			"email": "",
		}
	}
	data["author"] = author
	return data, nil
}

type fmInstall struct {
	Title                string `json:"title" binding:"required"`
	Subhead              string `json:"subhead" binding:"required"`
	Name                 string `json:"name" binding:"required"`
	Email                string `json:"email" binding:"email"`
	Password             string `json:"password" binding:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"eqfield=Password"`
}

func (p *HomePlugin) postInstall(l string, c *gin.Context) (interface{}, error) {
	var fm fmInstall
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}

	now := time.Now()
	ip := c.ClientIP()

	if err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		cnt, err := tx.Model(&User{}).Count()
		if err != nil {
			return err
		}
		if cnt > 0 {
			return p.I18n.E(l, "errors.forbidden")
		}
		if err = p.I18n.Set(tx, l, "site.title", fm.Title); err != nil {
			return err
		}
		if err = p.I18n.Set(tx, l, "site.subhead", fm.Subhead); err != nil {
			return err
		}
		user, err := p.Dao.AddEmailUser(tx, fm.Name, fm.Email, fm.Password)
		if err != nil {
			return err
		}
		if err = p.Dao.AddLog(tx, user.ID, ip, l, "nut.logs.sign-up"); err != nil {
			return err
		}
		user.ConfirmedAt = &now
		user.UpdatedAt = now
		if _, err = tx.Model(user).Column("confirmed_at", "updated_at").Update(); err != nil {
			return err
		}
		if err = p.Dao.AddLog(tx, user.ID, ip, l, "nut.logs.confirm"); err != nil {
			return err
		}
		for _, rn := range []string{RoleRoot, RoleAdmin} {
			if err := p.Dao.Allow(tx, user.ID, rn, DefaultResourceType, DefaultResourceID, 50, 0, 0); err != nil {
				return err
			}
			if err := p.Dao.AddLog(tx, user.ID, ip, l, "nut.logs.apply-role", rn, DefaultResourceType, DefaultResourceID); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}

type fmLeaveWord struct {
	Body string `json:"body" binding:"required"`
	Type string `json:"type" binding:"required"`
}

func (p *HomePlugin) createLeaveWord(l string, c *gin.Context) (interface{}, error) {
	var fm fmLeaveWord
	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	if err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		return tx.Insert(&LeaveWord{
			Body: fm.Body,
			Type: fm.Type,
		})
	}); err != nil {
		return nil, err
	}
	return gin.H{}, nil
}
