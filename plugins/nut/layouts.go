package nut

import (
	"net/http"

	"github.com/kapmahc/axe/web"
)

const (
	// NOTICE notice
	NOTICE = "notice"
	// WARNING warning
	WARNING = "warning"
	// ERROR error
	ERROR = "error"
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
			ss.AddFlash(err.Error(), ERROR)
			ctx.Save(ss)
			ctx.Redirect(http.StatusFound, fto)
		}
	}
}

// Application application layout
func Application(tpl string, fn func(string, web.H, *web.Context) error) http.HandlerFunc {
	return renderLayout("layouts/application/", tpl, fn)
}

// Dashboard dashboard layout
func Dashboard(tpl string, fn func(string, web.H, *web.Context) error) http.HandlerFunc {
	return renderLayout("layouts/dashboard/", tpl, fn)
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

		data := web.H{
			"locale":    lang,
			"languages": I18N().Languages(),
			"flashes":   flashes,
		}

		if err := fn(lang, data, ctx); err == nil {
			ctx.HTML(http.StatusOK, lyt+"index", tpl, data)
		} else {
			ctx.Error(http.StatusInternalServerError, lyt+"index", lyt+"error", err)
		}
	}
}
