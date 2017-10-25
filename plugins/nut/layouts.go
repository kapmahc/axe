package nut

import (
	"net/http"

	"github.com/kapmahc/axe/web"
	log "github.com/sirupsen/logrus"
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
)

// Layout layout
type Layout struct {
	Wrapper  *web.Wrapper  `inject:""`
	Settings *web.Settings `inject:""`
	I18n     *web.I18n     `inject:""`
}

// Redirect redirect
func (p *Layout) Redirect(to string, fn func(string, *web.Context) error) http.HandlerFunc {
	return p.Wrapper.HTTP(func(ctx *web.Context) {
		lang := ctx.Get(web.LOCALE).(string)
		if err := fn(lang, ctx); err != nil {
			log.Error(err)
			ss := ctx.Session()
			ss.AddFlash(err.Error(), ERROR)
			ctx.Save(ss)
		}
		ctx.Redirect(http.StatusFound, to)
	})
}

// JSON render json
func (p *Layout) JSON(fn func(string, *web.Context) (interface{}, error)) http.HandlerFunc {
	return p.Wrapper.HTTP(func(ctx *web.Context) {
		lang := ctx.Get(web.LOCALE).(string)
		if val, err := fn(lang, ctx); err == nil {
			ctx.JSON(http.StatusOK, val)
		} else {
			ctx.Abort(http.StatusInternalServerError, err)
		}
	})
}

// XML render xml
func (p *Layout) XML(fn func(string, *web.Context) (interface{}, error)) http.HandlerFunc {
	return p.Wrapper.HTTP(func(ctx *web.Context) {
		lang := ctx.Get(web.LOCALE).(string)
		if val, err := fn(lang, ctx); err == nil {
			ctx.XML(http.StatusOK, val)
		} else {
			log.Error(err)
			ctx.Abort(http.StatusInternalServerError, err)
		}
	})
}

// Application application layout
func (p *Layout) Application(tpl string, fn func(string, web.H, *web.Context) error) http.HandlerFunc {
	return p.renderLayout("layouts/application/index", tpl, fn)
}

func (p *Layout) renderLayout(lyt, tpl string, fn func(string, web.H, *web.Context) error) http.HandlerFunc {
	return p.Wrapper.HTTP(func(ctx *web.Context) {
		flashes := web.H{}
		ss := ctx.Session()
		for _, n := range []string{NOTICE, WARNING, ERROR} {
			flashes[n] = ss.Flashes(n)
		}
		ctx.Save(ss)
		lang := ctx.Get(web.LOCALE).(string)
		data := web.H{}

		if err := fn(lang, data, ctx); err != nil {
			ctx.Error(http.StatusInternalServerError, "layouts/application/index", "nut/error", err)
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

		data["locale"] = lang
		data["favicon"] = favicon
		data["languages"] = langs
		data["flashes"] = flashes

		// log.Debugf("%+v", data)
		ctx.HTML(http.StatusOK, lyt, tpl, data)
	})
}
