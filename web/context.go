package web

import (
	"encoding/json"
	"html/template"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/form"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	validator "gopkg.in/go-playground/validator.v9"
)

// Wrap wrap handler func
func Wrap(hf HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hf(NewContext(w, r))
	}
}

// HandlerFunc handler func
type HandlerFunc func(*Context)

// SetContext set context
func SetContext(secret []byte, views string, fm template.FuncMap, debug bool) {
	_sessionStore = sessions.NewCookieStore(secret)
	_render = render.New(render.Options{
		Directory:     views,
		Extensions:    []string{".html"},
		Funcs:         []template.FuncMap{fm},
		IsDevelopment: debug,
	})
}

// NewContext create a context
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
	}
}

var (
	_decoder      = form.NewDecoder()
	_validate     = validator.New()
	_render       *render.Render
	_sessionStore sessions.Store
)

// Context http context
type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
}

// Abort render error text
func (p *Context) Abort(s int, e error) {
	_render.Text(p.Writer, s, e.Error())
}

// Error render error page
func (p *Context) Error(s int, l, t string, e error) {
	p.HTML(
		http.StatusOK,
		l,
		t,
		H{"error": e.Error(), "created": time.Now(), "code": s},
	)
}

// HTML render html
func (p *Context) HTML(s int, l, t string, v interface{}) {
	_render.HTML(p.Writer, http.StatusOK, t, v, render.HTMLOptions{Layout: l})
}

// Text render text
func (p *Context) Text(s int, v string) {
	_render.Text(p.Writer, http.StatusOK, v)
}

// JSON render json
func (p *Context) JSON(s int, v interface{}) {
	_render.JSON(p.Writer, http.StatusOK, v)
}

// XML render xml
func (p *Context) XML(s int, v interface{}) {
	_render.XML(p.Writer, http.StatusOK, v)
}

// Session get session
func (p *Context) Session() *sessions.Session {
	ss, er := _sessionStore.Get(p.Request, "_session_")
	if er != nil {
		log.Error(er)
	}
	ss.Options = &sessions.Options{
		MaxAge:   0,
		HttpOnly: true,
		Secure:   p.Request.TLS != nil,
	}
	return ss
}

// Save save session
func (p *Context) Save(ss *sessions.Session) {
	if err := ss.Save(p.Request, p.Writer); err != nil {
		log.Error(err)
	}
}

// Get get value
func (p *Context) Get(k string) interface{} {
	return p.Request.Context().Value(K(k))
}

// Header get header
func (p *Context) Header(k string) string {
	return p.Request.Header.Get(k)
}

// ClientIP http client ip
func (p *Context) ClientIP() string {

	ip := p.Header("X-Forwarded-For")
	if idx := strings.IndexByte(ip, ','); idx >= 0 {
		ip = strings.TrimSpace(ip[0:idx])
	}
	if ip != "" {
		return ip
	}

	if ip = strings.TrimSpace(p.Header("X-Real-Ip")); ip != "" {
		return ip
	}

	if ip := p.Header("X-Appengine-Remote-Addr"); ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(p.Request.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// Bind bind http request json body and validate
func (p *Context) Bind(fm interface{}) error {
	dec := json.NewDecoder(p.Request.Body)
	if err := dec.Decode(fm); err != nil {
		return err
	}
	return _validate.Struct(fm)
}

// Form bind http form and validate
func (p *Context) Form(fm interface{}) error {
	if err := p.Request.ParseForm(); err != nil {
		return err
	}
	if err := _decoder.Decode(fm, p.Request.Form); err != nil {
		return err
	}
	return _validate.Struct(fm)
}

// Redirect redirect
func (p *Context) Redirect(s int, u string) {
	http.Redirect(p.Writer, p.Request, u, s)
}
