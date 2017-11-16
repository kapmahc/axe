package web

import (
	"net/http"

	"github.com/unrolled/render"
)

// H hash
type H map[string]interface{}

// HandlerFunc http handler func
type HandlerFunc func(*Context)

// Context http context
type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
	render  *render.Render
}

// Query get url query
func (p *Context) Query(k string) string {
	return p.Request.URL.Query().Get(k)
}

// JSON render json
func (p *Context) JSON(status int, value interface{}) {
	p.render.JSON(p.Writer, status, value)
}

// XML render xml
func (p *Context) XML(status int, value interface{}) {
	p.render.XML(p.Writer, status, value)
}

// HTML render html
func (p *Context) HTML(status int, layout, name string, value interface{}) {
	p.render.HTML(p.Writer, status, name, value, render.HTMLOptions{Layout: layout})
}
