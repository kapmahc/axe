package nut

import (
	"fmt"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
	log "github.com/sirupsen/logrus"
)

func (p *UsersPlugin) deleteSignOut(l string, c *web.Context) (interface{}, error) {
	user, err := p.Layout.CurrentUser(c)
	if err != nil {
		return nil, err
	}
	if err = p.DB.RunInTransaction(func(tx *pg.Tx) error {
		return p.Dao.AddLog(tx, user.ID, c.ClientIP(), l, "nut.logs.sign-out")
	}); err != nil {
		return nil, err
	}
	return web.H{}, nil
}

func (p *UsersPlugin) getChangePassword(l string, c *web.Context) (web.H, error) {
	return web.H{web.TITLE: p.I18n.T(l, "nut.users.change-password.title")}, nil
}

type fmUsersChangePassword struct {
	CurrentPassword      string `form:"currentPassword" validate:"required"`
	NewPassword          string `form:"newPassword" validate:"required,min=6"`
	PasswordConfirmation string `form:"passwordConfirmation" validate:"eqfield=NewPassword"`
}

func (p *UsersPlugin) postChangePassword(l string, c *web.Context) (interface{}, error) {
	var fm fmUsersChangePassword
	if err := c.Bind(&fm); err != nil {
		return nil, err
	}
	user, err := p.Layout.CurrentUser(c)
	if !p.Security.Check(user.Password, []byte(fm.CurrentPassword)) {
		return nil, p.I18n.E(l, "nut.errors.user-bad-password")
	}

	if user.Password, err = p.Security.Hash([]byte(fm.NewPassword)); err != nil {
		return nil, err
	}
	user.UpdatedAt = time.Now()
	err = p.DB.RunInTransaction(func(tx *pg.Tx) error {
		if _, er := tx.Model(user).Column("name", "password").Update(); er != nil {
			return err
		}
		return p.Dao.AddLog(tx, user.ID, c.ClientIP(), l, "nut.logs.change-password")
	})
	return web.H{}, err
}

func (p *UsersPlugin) getProfile(l string, c *web.Context) (web.H, error) {
	return web.H{web.TITLE: p.I18n.T(l, "nut.users.profile.title")}, nil
}

type fmUsersProfile struct {
	Name string `form:"name" validate:"required"`
}

func (p *UsersPlugin) postProfile(l string, c *web.Context) (interface{}, error) {
	var fm fmUsersProfile
	if err := c.Bind(&fm); err != nil {
		return nil, err
	}
	user, err := p.Layout.CurrentUser(c)
	if err != nil {
		return nil, err
	}
	user.Name = fm.Name
	user.UpdatedAt = time.Now()
	_, err = p.DB.Model(user).Column("name", "updated_at").Update()
	return web.H{}, err
}

func (p *UsersPlugin) getLogs(l string, c *web.Context) (web.H, error) {
	user, err := p.Layout.CurrentUser(c)
	if err != nil {
		return nil, err
	}
	var items []Log
	if err := p.DB.Model(&items).
		Column("id", "ip", "message", "created_at").
		Where("user_id = ?", user.ID).
		Order("created_at DESC").
		Select(); err != nil {
		return nil, err
	}
	return web.H{"items": items}, nil
}

func (p *UsersPlugin) getSignIn(l string, c *web.Context) (web.H, error) {
	return web.H{web.TITLE: p.I18n.T(l, "nut.users.sign-in.title")}, nil
}

type fmUsersSignIn struct {
	Email    string `form:"email" validate:"email"`
	Password string `form:"password" validate:"required"`
}

func (p *UsersPlugin) postSignIn(l string, c *web.Context) (interface{}, error) {
	var fm fmUsersSignIn
	if err := c.Bind(&fm); err != nil {
		return nil, err
	}

	ss := c.Session()
	if err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		user, err := p.Dao.SignIn(tx, l, c.ClientIP(), fm.Email, fm.Password)
		if err != nil {
			return err
		}
		ss.Values["currentUser"] = web.H{"uid": user.UID, "name": user.Name}
		return nil
	}); err != nil {
		return nil, err
	}
	c.Save(ss)
	return web.H{}, nil
}

func (p *UsersPlugin) getSignUp(l string, c *web.Context) (web.H, error) {
	return web.H{web.TITLE: p.I18n.T(l, "nut.users.sign-up.title")}, nil
}

type fmUsersSignUp struct {
	Name                 string `form:"name" validate:"required"`
	Email                string `form:"email" validate:"email"`
	Password             string `form:"password" validate:"required,min=6"`
	PasswordConfirmation string `form:"passwordConfirmation" validate:"eqfield=Password"`
}

func (p *UsersPlugin) postSignUp(l string, c *web.Context) (interface{}, error) {
	var fm fmUsersSignUp
	if err := c.Bind(&fm); err != nil {
		return nil, err
	}

	ip := c.ClientIP()
	if err := p.DB.RunInTransaction(func(tx *pg.Tx) error {
		user, err := p.Dao.AddEmailUser(tx, fm.Name, fm.Email, fm.Password)
		if err != nil {
			return err
		}
		if err = p.Dao.AddLog(tx, user.ID, ip, l, "nut.logs.sign-up"); err != nil {
			return err
		}
		if err = p.sendEmail(c, l, user, actConfirm); err != nil {
			log.Error(err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return web.H{web.MESSAGE: p.I18n.T(l, "nut.users.confirm.notice")}, nil
}

type fmUsersEmail struct {
	Email string `form:"email" validate:"email"`
}

func (p *UsersPlugin) getEmailForm(act string) web.HTMLHandlerFunc {
	return func(l string, c *web.Context) (web.H, error) {
		return web.H{
			"action":  act,
			web.TITLE: p.I18n.T(l, "nut.users."+act+".title"),
		}, nil
	}
}

func (p *UsersPlugin) getConfirmToken(l string, c *web.Context) error {
	cm, err := p.Jwt.Validate([]byte(c.Params["token"]))
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
	if err = p.DB.RunInTransaction(func(tx *pg.Tx) error {
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
	ss.AddFlash(p.I18n.T(l, "nut.users.confirm.success"), web.NOTICE)
	c.Save(ss)

	return nil
}

func (p *UsersPlugin) postConfirm(l string, c *web.Context) (interface{}, error) {
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
	if err := p.sendEmail(c, l, &user, actConfirm); err != nil {
		log.Error(err)
	}

	return web.H{web.MESSAGE: p.I18n.T(l, "nut.users.confirm.notice")}, nil
}

func (p *UsersPlugin) getUnlockToken(l string, c *web.Context) error {
	cm, err := p.Jwt.Validate([]byte(c.Params["token"]))
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
	if err = p.DB.RunInTransaction(func(tx *pg.Tx) error {
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
	ss.AddFlash(p.I18n.T(l, "nut.users.unlock.success"), web.NOTICE)
	c.Save(ss)
	return nil
}

func (p *UsersPlugin) postUnlock(l string, c *web.Context) (interface{}, error) {
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
	if err := p.sendEmail(c, l, &user, actUnlock); err != nil {
		log.Error(err)
	}

	return web.H{web.MESSAGE: p.I18n.T(l, "nut.users.unlock.notice")}, nil
}

func (p *UsersPlugin) postForgotPassword(l string, c *web.Context) (interface{}, error) {
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

	if err := p.sendEmail(c, l, &user, actResetPassword); err != nil {
		log.Error(err)
	}

	return web.H{web.MESSAGE: p.I18n.T(l, "nut.users.forgot-password.notice")}, nil
}

func (p *UsersPlugin) getResetPassword(l string, c *web.Context) (web.H, error) {
	return web.H{web.TITLE: p.I18n.T(l, "nut.users.reset-password.title")}, nil
}

type fmUsersResetPassword struct {
	Password             string `form:"password" validate:"required"`
	PasswordConfirmation string `form:"passwordConfirmation" validate:"eqfield=Password"`
}

func (p *UsersPlugin) postResetPassword(l string, c *web.Context) (interface{}, error) {
	var fm fmUsersResetPassword
	if err := c.Bind(&fm); err != nil {
		return nil, err
	}

	cm, err := p.Jwt.Validate([]byte(c.Params["token"]))
	if err != nil {
		return nil, err
	}
	if cm.Get("act").(string) != actResetPassword {
		return nil, p.I18n.E(l, "errors.bad-action")
	}
	var user User
	if err = p.DB.Model(&user).
		Column("id", "locked_at").
		Where("uid = ?", cm.Get("uid")).
		Limit(1).Select(); err != nil {
		return nil, err
	}
	now := time.Now()
	if user.Password, err = p.Security.Hash([]byte(fm.Password)); err != nil {
		return nil, err
	}
	user.UpdatedAt = now
	if err = p.DB.RunInTransaction(func(tx *pg.Tx) error {
		if _, err = tx.Model(&user).Column("password", "updated_at").Update(); err != nil {
			return err
		}
		if err = p.Dao.AddLog(tx, user.ID, c.ClientIP(), l, "nut.logs.reset-password"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return web.H{web.MESSAGE: p.I18n.T(l, "nut.users.reset-password.success")}, nil
}

const (
	actConfirm       = "confirm"
	actUnlock        = "unlock"
	actResetPassword = "reset-password"

	// SendEmailJob send email
	SendEmailJob = "send.email"
)

func (p *UsersPlugin) sendEmail(c *web.Context, lang string, user *User, act string) error {

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
		Home:  c.Home(),
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
