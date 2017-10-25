package nut

import (
	"fmt"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
	log "github.com/sirupsen/logrus"
)

type fmUsersSignIn struct {
}

func (p *UsersPlugin) postUsersSignIn(l string, c *web.Context) (interface{}, error) {
	return web.H{}, nil
}

type fmUsersSignUp struct {
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"email"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"eqfield=Password"`
}

func (p *UsersPlugin) postUsersSignUp(l string, c *web.Context) (interface{}, error) {
	var fm fmUsersSignUp
	if err := c.Bind(&fm); err != nil {
		return nil, err
	}

	ip := c.ClientIP()
	if err := web.Tx(p.DB, func(tx *pg.Tx) error {
		user, err := p.Dao.AddEmailUser(tx, fm.Name, fm.Email, fm.Password)
		if err != nil {
			return err
		}
		if err = p.Dao.AddLog(tx, user.ID, ip, l, "nut.logs.sign-up"); err != nil {
			return err
		}
		if err = p.sendEmail(c, user, actConfirm); err != nil {
			log.Error(err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return web.H{}, nil
}

type fmUsersEmail struct {
	Email string `json:"email" validate:"email"`
}

func (p *UsersPlugin) getUsersConfirmToken(l string, c *web.Context) error {
	cm, err := p.Jwt.Validate([]byte(c.Param("token")))
	if err != nil {
		return err
	}
	if cm.Get("act").(string) != actConfirm {
		return p.I18n.E(l, "errors.bad-action")
	}
	var user User
	if err = p.DB.Model(&user).
		Column("id", "confirmed_at").
		Where("uid = ?", cm.Get("uid")).
		Limit(1).Select(); err != nil {
		return err
	}
	if user.IsConfirm() {
		return p.I18n.E(l, "nut.errors.user-already-confirm")
	}

	now := time.Now()
	if err = web.Tx(p.DB, func(tx *pg.Tx) error {
		user.ConfirmedAt = &now
		user.UpdatedAt = now
		if _, err = tx.Model(&user).Column("confirmed_at", "updated_at").Update(); err != nil {
			return err
		}
		if err = p.Dao.AddLog(tx, user.ID, c.ClientIP(), l, "nut.logs.confirm"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	ss := c.Session()
	ss.AddFlash(p.I18n.T(l, "nut.users.confirm.success"), NOTICE)
	c.Save(ss)

	return nil
}

func (p *UsersPlugin) postUsersConfirm(l string, c *web.Context) (interface{}, error) {
	var fm fmUsersEmail
	if err := c.Bind(&fm); err != nil {
		return nil, err
	}

	var user User
	if err := p.DB.Model(&user).
		Column("uid", "confirmed_at", "email").
		Where("provider_type = ? AND provider_id = ?", UserTypeEmail, fm.Email).
		Limit(1).Select(); err != nil {
		return nil, err
	}
	if user.IsConfirm() {
		return nil, p.I18n.E(l, "nut.errors.user-already-confirm")
	}
	if err := p.sendEmail(c, &user, actConfirm); err != nil {
		log.Error(err)
	}

	return web.H{}, nil
}

func (p *UsersPlugin) getUsersUnlockToken(l string, c *web.Context) error {
	cm, err := p.Jwt.Validate([]byte(c.Param("token")))
	if err != nil {
		return err
	}
	if cm.Get("act").(string) != actUnlock {
		return p.I18n.E(l, "errors.bad-action")
	}
	var user User
	if err = p.DB.Model(&user).
		Column("id", "locked_at").
		Where("uid = ?", cm.Get("uid")).
		Limit(1).Select(); err != nil {
		return err
	}
	if !user.IsLock() {
		return p.I18n.E(l, "nut.errors.user-not-lock")
	}

	now := time.Now()
	if err = web.Tx(p.DB, func(tx *pg.Tx) error {
		user.LockedAt = nil
		user.UpdatedAt = now
		if _, err = tx.Model(&user).Column("locked_at", "updated_at").Update(); err != nil {
			return err
		}
		if err = p.Dao.AddLog(tx, user.ID, c.ClientIP(), l, "nut.logs.unlock"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	ss := c.Session()
	ss.AddFlash(p.I18n.T(l, "nut.users.unlock.success"), NOTICE)
	c.Save(ss)

	return nil
}

func (p *UsersPlugin) postUsersUnlock(l string, c *web.Context) (interface{}, error) {
	var fm fmUsersEmail
	if err := c.Bind(&fm); err != nil {
		return nil, err
	}

	var user User
	if err := p.DB.Model(&user).
		Column("uid", "locked_at", "email").
		Where("provider_type = ? AND provider_id = ?", UserTypeEmail, fm.Email).
		Limit(1).Select(); err != nil {
		return nil, err
	}
	if !user.IsLock() {
		return nil, p.I18n.E(l, "nut.errors.user-not-lock")
	}
	if err := p.sendEmail(c, &user, actUnlock); err != nil {
		log.Error(err)
	}

	return web.H{}, nil
}

func (p *UsersPlugin) postUsersForgotPassword(l string, c *web.Context) (interface{}, error) {
	var fm fmUsersEmail
	if err := c.Bind(&fm); err != nil {
		return nil, err
	}

	var user User
	if err := p.DB.Model(&user).
		Column("uid", "email").
		Where("provider_type = ? AND provider_id = ?", UserTypeEmail, fm.Email).
		Limit(1).Select(); err != nil {
		return nil, err
	}

	if err := p.sendEmail(c, &user, actResetPassword); err != nil {
		log.Error(err)
	}

	return web.H{}, nil
}

const (
	actConfirm       = "confirm"
	actUnlock        = "unlock"
	actResetPassword = "reset-password"

	// SendEmailJob send email
	SendEmailJob = "send.email"
)

func (p *UsersPlugin) sendEmail(ctx *web.Context, user *User, act string) error {
	lang := ctx.Get(web.LOCALE).(string)
	cm := jws.Claims{}
	cm.Set("act", act)
	cm.Set("uid", user.UID)
	tkn, err := p.Jwt.Sum(cm, time.Hour*6)
	if err != nil {
		return err
	}

	obj := struct {
		Home  string
		Token string
	}{
		Home:  ctx.Home(),
		Token: string(tkn),
	}

	subject, err := p.I18n.H(lang, fmt.Sprintf("nut.emails.%s.subject", act), obj)
	if err != nil {
		return err
	}
	body, err := p.I18n.H(lang, fmt.Sprintf("nut.emails.%s.body", act), obj)
	if err != nil {
		return err
	}

	return p.Jobber.Send(SendEmailJob, 0, map[string]string{
		"to":      user.Email,
		"subject": subject,
		"body":    body,
	})

}
