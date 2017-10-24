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

// Form form handle
// func Form(sto, fto string, fm interface{}, fn func(string, interface{}, *web.Context) error) http.HandlerFunc {
// 	return func(wrt http.ResponseWriter, req *http.Request) {
// 		ctx := web.NewContext(wrt, req)
// 		lang := ctx.Get(web.LOCALE).(string)
// 		err := ctx.Bind(fm)
// 		if err == nil {
// 			err = fn(lang, fm, ctx)
// 		}
// 		if err == nil {
// 			ctx.Redirect(http.StatusFound, sto)
// 		} else {
// 			ss := ctx.Session()
// 			if ve, ok := err.(validator.ValidationErrors); ok {
// 				for _, er := range ve {
// 					ss.AddFlash(fmt.Sprintf("Validation for '%s' failed on the '%s' tag", er.Field(), er.Tag()), ERROR)
// 				}
// 			} else {
// 				ss.AddFlash(err.Error(), ERROR)
// 			}
// 			ctx.Save(ss)
// 			ctx.Redirect(http.StatusFound, fto)
// 		}
// 	}
// }

// Application application layout
func (p *Layout) Application(tpl string, fn func(string, web.H, *web.Context) error) http.HandlerFunc {
	return p.renderLayout("layouts/application/index", tpl, fn)
}

func (p *Layout) renderLayout(lyt, tpl string, fn func(string, web.H, *web.Context) error) http.HandlerFunc {
	return p.Wrapper.HTTP(func(ctx *web.Context) {
		lang := ctx.Get(web.LOCALE).(string)

		flashes := web.H{}
		ss := ctx.Session()
		for _, n := range []string{NOTICE, WARNING, ERROR} {
			flashes[n] = ss.Flashes(n)
		}
		ctx.Save(ss)

		var favicon string
		if err := p.Settings.Get("site.favicon", &favicon); err != nil {
			favicon = "/assets/favicon.png"
		}
		var author map[string]interface{}
		if err := p.Settings.Get("site.author", &author); err != nil {
			author = web.H{}
		}
		langs, err := p.I18n.Languages()
		if err != nil {
			langs = make([]string, 0)
		}

		data := web.H{
			"locale":    lang,
			"favicon":   favicon,
			"author":    author,
			"languages": langs,
			"flashes":   flashes,
			// csrf.TemplateTag: csrf.TemplateField(req),
			// "_csrf_token":    csrf.Token(req),
		}

		if err := fn(lang, data, ctx); err == nil {
			ctx.HTML(http.StatusOK, lyt, tpl, data)
		} else {
			ctx.Error(http.StatusInternalServerError, "layouts/application/index", "nut/error", err)
		}
	})
}
