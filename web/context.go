package web

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-playground/form"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	validator "gopkg.in/go-playground/validator.v9"
)

// Context http context
type Context struct {
	Request  *http.Request
	Writer   http.ResponseWriter
	render   *render.Render
	store    sessions.Store
	decoder  *form.Decoder
	validate *validator.Validate
	query    url.Values
	params   map[string]string
}

// Param http request url param
func (p *Context) Param(key string) string {
	return p.params[key]
}

// Query http request query param
func (p *Context) Query(key string) string {
	return p.query.Get(key)
}

// Home home url
func (p *Context) Home() string {
	scheme := "http"
	if p.Request.TLS != nil {
		scheme += "s"
	}
	return scheme + "://" + p.Request.Host
}

// Abort render error text
func (p *Context) Abort(s int, e error) {
	log.Error(e)
	p.render.Text(p.Writer, s, e.Error())
}

// Error render error page
func (p *Context) Error(s int, l, t string, e error) {
	log.Error(e)
	p.HTML(
		http.StatusOK,
		l,
		t,
		H{"error": e.Error(), "created": time.Now(), "code": s},
	)
}

// HTML render html
func (p *Context) HTML(s int, l, t string, v interface{}) {
	p.render.HTML(p.Writer, http.StatusOK, t, v, render.HTMLOptions{Layout: l})
}

// Text render text
func (p *Context) Text(s int, v string) {
	p.render.Text(p.Writer, http.StatusOK, v)
}

// JSON render json
func (p *Context) JSON(s int, v interface{}) {
	p.render.JSON(p.Writer, http.StatusOK, v)
}

// XML render xml
func (p *Context) XML(s int, v interface{}) {
	p.render.XML(p.Writer, http.StatusOK, v)
}

// Session get session
func (p *Context) Session() *sessions.Session {
	ss, er := p.store.Get(p.Request, "_session_")
	if er != nil {
		log.Error(er)
	}
	ss.Options = &sessions.Options{
		Path:     "/",
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
	return p.validate.Struct(fm)
}

// Form bind http form and validate
func (p *Context) Form(fm interface{}) error {
	if err := p.Request.ParseForm(); err != nil {
		return err
	}
	if err := p.decoder.Decode(fm, p.Request.Form); err != nil {
		return err
	}
	return p.validate.Struct(fm)
}

// Redirect redirect
func (p *Context) Redirect(s int, u string) {
	http.Redirect(p.Writer, p.Request, u, s)
}
