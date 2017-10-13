package nut

import (
	"fmt"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/kapmahc/axe/web"
	validator "gopkg.in/go-playground/validator.v9"
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

// JSON render json
func JSON(fn func(string, *web.Context) (interface{}, error)) http.HandlerFunc {
	return func(wrt http.ResponseWriter, req *http.Request) {
		ctx := web.NewContext(wrt, req)
		lang := ctx.Get(web.LOCALE).(string)
		if val, err := fn(lang, ctx); err == nil {
			ctx.JSON(http.StatusOK, val)
		} else {
			ctx.Abort(http.StatusInternalServerError, err)
		}
	}
}

// XML render xml
func XML(fn func(string, *web.Context) (interface{}, error)) http.HandlerFunc {
	return func(wrt http.ResponseWriter, req *http.Request) {
		ctx := web.NewContext(wrt, req)
		lang := ctx.Get(web.LOCALE).(string)
		if val, err := fn(lang, ctx); err == nil {
			ctx.XML(http.StatusOK, val)
		} else {
			ctx.Abort(http.StatusInternalServerError, err)
		}
	}
}

// Form form handle
func Form(sto, fto string, fm interface{}, fn func(string, interface{}, *web.Context) error) http.HandlerFunc {
	return func(wrt http.ResponseWriter, req *http.Request) {
		ctx := web.NewContext(wrt, req)
		lang := ctx.Get(web.LOCALE).(string)
		err := ctx.Bind(fm)
		if err == nil {
			err = fn(lang, fm, ctx)
		}
		if err == nil {
			ctx.Redirect(http.StatusFound, sto)
		} else {
			ss := ctx.Session()
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, er := range ve {
					ss.AddFlash(fmt.Sprintf("Validation for '%s' failed on the '%s' tag", er.Field(), er.Tag()), ERROR)
				}
			} else {
				ss.AddFlash(err.Error(), ERROR)
			}
			ctx.Save(ss)
			ctx.Redirect(http.StatusFound, fto)
		}
	}
}

// Application application layout
func Application(tpl string, fn func(string, web.H, *web.Context) error) http.HandlerFunc {
	return renderLayout("layouts/application/index", tpl, fn)
}

// Dashboard dashboard layout
func Dashboard(tpl string, fn func(string, web.H, *web.Context) error) http.HandlerFunc {
	return renderLayout("layouts/dashboard/index", tpl, fn)
}

func renderLayout(lyt, tpl string, fn func(string, web.H, *web.Context) error) http.HandlerFunc {
	return func(wrt http.ResponseWriter, req *http.Request) {
		ctx := web.NewContext(wrt, req)
		lang := ctx.Get(web.LOCALE).(string)

		flashes := web.H{}
		ss := ctx.Session()
		for _, n := range []string{NOTICE, WARNING, ERROR} {
			flashes[n] = ss.Flashes(n)
		}
		ctx.Save(ss)

		var favicon string
		if err := SETTINGS().Get("site.favicon", &favicon); err != nil {
			favicon = "/assets/favicon.png"
		}
		var author map[string]interface{}
		if err := SETTINGS().Get("site.author", &author); err != nil {
			author = web.H{}
		}

		data := web.H{
			"locale":         lang,
			"favicon":        favicon,
			"author":         author,
			"languages":      I18N().Languages(),
			"flashes":        flashes,
			csrf.TemplateTag: csrf.TemplateField(req),
			"_csrf_token":    csrf.Token(req),
		}

		if err := fn(lang, data, ctx); err == nil {
			ctx.HTML(http.StatusOK, lyt, tpl, data)
		} else {
			ctx.Error(http.StatusInternalServerError, "layouts/application/index", "nut/error", err)
		}
	}
}
