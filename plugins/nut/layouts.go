package nut

import (
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
)

// Layout layout
type Layout struct {
	Settings *web.Settings `inject:""`
	Router   *web.Router   `inject:""`
	I18n     *web.I18n     `inject:""`
	Jwt      *web.Jwt      `inject:""`
	DB       *pg.DB        `inject:""`
	Dao      *Dao          `inject:""`
}

type fmUEditor struct {
	Body string `form:"body" binding:"required"`
	Next string `form:"next" binding:"required"`
}

func (p *Layout) checkToken(act string, c *web.Context, check func(*User, uint) bool) (string, uint, error) {
	lng := c.Locale()
	token := c.Params["token"]
	cw, err := p.Jwt.Validate([]byte(token))
	if err != nil {
		return "", 0, err
	}
	if cw.Get("act").(string) != act {
		return "", 0, p.I18n.E(lng, "errors.bad-action")
	}
	var user User
	if err := p.DB.Model(&user).Where("uid = ?", cw.Get("uid")).Limit(1).Select(); err != nil {
		return "", 0, p.I18n.E(lng, "errors.forbidden")
	}
	tid := uint(cw.Get("tid").(float64))
	if user.IsLock() || !user.IsConfirm() || !check(&user, tid) {
		return "", 0, p.I18n.E(lng, "errors.forbidden")
	}

	return token, tid, nil
}

// CurrentUser current user
func (p *Layout) CurrentUser(c *web.Context) (*User, error) {
	cm, err := p.Jwt.Parse(c.Request)
	if err == nil {
		var it User
		if err = p.DB.Model(&it).
			Where("uid = ?", cm.Get("uid")).
			Limit(1).Select(); err == nil && it.IsConfirm() && !it.IsLock() {
			return &it, nil
		}
	}
	return nil, p.I18n.E(c.Locale(), "errors.forbidden")
}
