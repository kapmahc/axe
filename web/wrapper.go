package web

import (
	"html/template"
	"net/http"

	"github.com/go-playground/form"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	validator "gopkg.in/go-playground/validator.v9"
)

// NewWrapper new wrapper
func NewWrapper(secret []byte, views string, fm template.FuncMap, debug bool) *Wrapper {
	return &Wrapper{
		decoder:  form.NewDecoder(),
		validate: validator.New(),
		store:    sessions.NewCookieStore(secret),
		render: render.New(render.Options{
			Directory:     views,
			Extensions:    []string{".html"},
			Funcs:         []template.FuncMap{fm},
			IsDevelopment: debug,
		}),
	}
}

// HandlerFunc handler func
type HandlerFunc func(*Context)

// Wrapper wrapper
type Wrapper struct {
	decoder  *form.Decoder
	validate *validator.Validate
	render   *render.Render
	store    sessions.Store
}

// Context create a context
func (p *Wrapper) Context(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:   w,
		Request:  r,
		render:   p.render,
		store:    p.store,
		decoder:  p.decoder,
		validate: p.validate,
	}
}

// HTTP http handler
func (p *Wrapper) HTTP(hf HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hf(p.Context(w, r))
	}
}
