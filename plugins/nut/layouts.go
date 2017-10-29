package nut

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
	log "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v8"
)

const (
	// NOTICE notice
	NOTICE = "notice"
	// WARNING warning
	WARNING = "warning"
	// ERROR error
	ERROR = "error"

	// TITLE title
	TITLE = "title"
	// CurrentUser current user
	CurrentUser = "currentUser"
	// IsAdmin is admin?
	IsAdmin = "isAdmin"
)

// Layout layout
type Layout struct {
	Settings *web.Settings `inject:""`
	I18n     *web.I18n     `inject:""`
	Jwt      *web.Jwt      `inject:""`
	DB       *pg.DB        `inject:""`
	Dao      *Dao          `inject:""`
}

// Home home url
func (p *Layout) Home(req *http.Request) string {
	scheme := "http"
	if req.TLS != nil {
		scheme += "s"
	}
	return scheme + "://" + req.Host
}

// CurrentUserMiddleware currend user middleware
func (p *Layout) CurrentUserMiddleware(c *gin.Context) {
	cm, err := p.Jwt.Parse(c.Request)
	if err == nil {
		var it User
		if err = p.DB.Model(&it).Where("uid = ?", cm.Get("uid")).Limit(1).Select(); err == nil {
			c.Set(CurrentUser, &it)
			c.Set(IsAdmin, p.Dao.Is(it.ID, RoleAdmin))
		}
	}
}

// MustSignInMiddleware currend user middleware
func (p *Layout) MustSignInMiddleware(c *gin.Context) {
	if it, ok := c.Get(CurrentUser); ok {
		if user, ok := it.(*User); ok && user.IsConfirm() && !user.IsLock() {
			return
		}
	}
	p.Abort(c, http.StatusForbidden, p.I18n.E(c.MustGet(web.LOCALE).(string), "errors.forbidden"))
}

// MustAdminMiddleware currend user middleware
func (p *Layout) MustAdminMiddleware(c *gin.Context) {
	if it, ok := c.Get(IsAdmin); ok {
		if is, ok := it.(bool); ok && is {
			return
		}
	}
	p.Abort(c, http.StatusForbidden, p.I18n.E(c.MustGet(web.LOCALE).(string), "errors.forbidden"))
}

// Redirect redirect
func (p *Layout) Redirect(to string, fn func(string, *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet(web.LOCALE).(string)
		if err := fn(l, c); err != nil {
			log.Error(err)
			ss := sessions.Default(c)
			ss.AddFlash(err.Error(), ERROR)
			ss.Save()
		}
		c.Redirect(http.StatusFound, to)
	}
}

// JSON render json
func (p *Layout) JSON(fn func(string, *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet(web.LOCALE).(string)
		if val, err := fn(l, c); err == nil {
			c.JSON(http.StatusOK, val)
		} else {
			p.Abort(c, http.StatusInternalServerError, err)
		}
	}
}

// Abort abort error
func (p *Layout) Abort(c *gin.Context, s int, e error) {
	log.Error(e)
	if er, ok := e.(validator.ValidationErrors); ok {
		var ss []string
		for _, it := range er {
			ss = append(ss, fmt.Sprintf("Validation for '%s' failed on the '%s' tag;", it.Field, it.Tag))
		}
		c.String(http.StatusBadRequest, strings.Join(ss, "\n"))
	} else {
		c.String(s, e.Error())
	}
	// c.AbortWithError(s, err)
}

// XML render xml
func (p *Layout) XML(fn func(string, *gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet(web.LOCALE).(string)
		if val, err := fn(l, c); err == nil {
			c.XML(http.StatusOK, val)
		} else {
			p.Abort(c, http.StatusInternalServerError, err)
		}
	}
}

// Application application layout
func (p *Layout) Application(tpl string, fn func(string, gin.H, *gin.Context) error) gin.HandlerFunc {
	return p.renderLayout(tpl, fn)
}

func (p *Layout) renderLayout(tpl string, fn func(string, gin.H, *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		flashes := gin.H{}
		ss := sessions.Default(c)
		for _, n := range []string{NOTICE, WARNING, ERROR} {
			flashes[n] = ss.Flashes(n)
		}
		ss.Save()
		lang := c.MustGet(web.LOCALE).(string)
		data := gin.H{}

		if err := fn(lang, data, c); err != nil {
			log.Error(err)
			c.HTML(http.StatusInternalServerError, "nut-error.html", gin.H{"error": err.Error()})
			return
		}

		var favicon string
		if err := p.Settings.Get("site.favicon", &favicon); err != nil {
			favicon = "/assets/favicon.png"
		}

		langs, err := p.I18n.Languages()
		if err != nil {
			langs = make([]string, 0)
		}

		// csrf.TemplateTag: csrf.TemplateField(req),
		// "_csrf_token":    csrf.Token(req),
		var author map[string]string
		if err := p.Settings.Get("site.author", &author); err != nil {
			author = make(map[string]string)
		}

		data["author"] = author
		data["locale"] = lang
		data["favicon"] = favicon
		data["languages"] = langs
		data["flashes"] = flashes

		// log.Debugf("%+v", data)
		c.HTML(http.StatusOK, tpl, data)
	}
}
