package nut

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
)

func (p *HomePlugin) getAPISiteInfo(l string, c *web.Context) (interface{}, error) {
	// -----------
	langs, err := p.I18n.Languages()
	if err != nil {
		return nil, err
	}
	data := web.H{"locale": l, "languages": langs}
	// -----------
	for _, k := range []string{"title", "subhead", "keywords", "description", "copyright"} {
		data[k] = p.I18n.T(l, "site."+k)
	}
	// -----------
	var author map[string]interface{}
	if err := p.Settings.Get("site.author", &author); err != nil {
		author = web.H{}
	}
	data["author"] = author
	return data, nil
}

type fmInstall struct {
	Title                string `json:"title" validate:"required"`
	Subhead              string `json:"subhead" validate:"required"`
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"email"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"eqfield=Password"`
}

func (p *HomePlugin) postInstall(l string, c *web.Context) (interface{}, error) {
	var fm fmInstall
	if err := c.Bind(&fm); err != nil {
		return nil, err
	}

	now := time.Now()
	ip := c.ClientIP()

	if err := web.Tx(p.DB, func(tx *pg.Tx) error {
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
	return web.H{}, nil
}
