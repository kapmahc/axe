package nut

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
)

func (p *HomePlugin) getHome(c *web.Context) {
	lang := c.Locale()
	theme := c.Query("theme")
	if theme == "" {
		if err := p.Settings.Get("site.home.theme", &theme); err != nil {
			theme = "off-canvas"
		}
	}
	web.HTML(web.APPLICATION, "nut/home/"+theme, func(_ string, _ *web.Context) (web.H, error) {
		return web.H{
			web.TITLE: p.I18n.T(lang, "nut.home.title"),
		}, nil
	})(c)
}

type fmInstall struct {
	Title                string `form:"title" validate:"required"`
	Subhead              string `form:"subhead" validate:"required"`
	Name                 string `form:"name" validate:"required"`
	Email                string `form:"email" validate:"email"`
	Password             string `form:"password" validate:"required"`
	PasswordConfirmation string `form:"passwordConfirmation" validate:"eqfield=Password"`
}

func (p *HomePlugin) getInstall(l string, c *web.Context) (web.H, error) {
	return web.H{web.TITLE: p.I18n.T(l, "nut.install.title")}, nil
}

func (p *HomePlugin) postInstall(l string, c *web.Context) (interface{}, error) {
	var fm fmInstall
	if err := c.Bind(&fm); err != nil {
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
	return web.H{}, nil
}

func (p *HomePlugin) newLeaveWord(l string, c *web.Context) (web.H, error) {
	return web.H{web.TITLE: p.I18n.T(l, "nut.leave-wods.new.title")}, nil
}

type fmLeaveWord struct {
	Body string `form:"body" validate:"required"`
	Type string `form:"type" validate:"required"`
}

func (p *HomePlugin) createLeaveWord(l string, c *web.Context) (interface{}, error) {
	var fm fmLeaveWord
	if err := c.Bind(&fm); err != nil {
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
	return web.H{}, nil
}
