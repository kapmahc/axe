package web

import (
	"math"
	"net"
	"net/http"
	"strings"

	"github.com/go-playground/form"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
	"golang.org/x/text/language"
	validator "gopkg.in/go-playground/validator.v9"
)

const (
	// NOTICE notice
	NOTICE = "notice"
	// WARNING warning
	WARNING = "warning"
	// ERROR = error
	ERROR = "error"
)

// NewContext create an context
func NewContext(req *http.Request, wrt http.ResponseWriter, ste sessions.Store, mch language.Matcher, rdr *render.Render, val *validator.Validate, dec *form.Decoder, langs ...language.Tag) *Context {
	return &Context{
		Writer:    wrt,
		Request:   req,
		Payload:   H{},
		matcher:   mch,
		validate:  val,
		decoder:   dec,
		render:    rdr,
		languages: langs,
	}
}

// H hash
type H map[string]interface{}

// HandlerFunc http handler func
type HandlerFunc func(*Context)

// Context http context
type Context struct {
	Request *http.Request
	Writer  http.ResponseWriter
	Payload H

	render   *render.Render
	matcher  language.Matcher
	validate *validator.Validate
	decoder  *form.Decoder
	store    sessions.Store

	locale    string
	languages []language.Tag
}

// Flashes read flashes
func (p *Context) Flashes() {
	flashes := H{}
	ss := p.Session()
	for _, k := range []string{NOTICE, WARNING, ERROR} {
		flashes[k] = ss.Flashes(k)
	}
	p.Save(ss)
	p.Payload["flashes"] = ss
}

// Session get session
func (p *Context) Session() *sessions.Session {
	s, _ := p.store.Get(p.Request, "session")
	return s
}

// Save save session
func (p *Context) Save(s *sessions.Session) error {
	s.Options.Path = "/"
	s.Options.MaxAge = 0
	s.Options.HttpOnly = true
	s.Options.Secure = p.Secure()
	return s.Save(p.Request, p.Writer)
}

// Header get http request header
func (p *Context) Header(k string) string {
	return p.Request.Header.Get(k)
}

// ClientIP client ip, orders: X-Forwarded-For, X-Real-Ip, RemoteAddr
func (p *Context) ClientIP() string {
	ip := p.Header("X-Forwarded-For")

	if idx := strings.IndexByte(ip, ','); idx >= 0 {
		ip = ip[0:idx]
	}
	ip = strings.TrimSpace(ip)
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(p.Header("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(p.Request.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// Bind bind form values
// https://github.com/go-playground/form
// https://github.com/go-playground/validator
func (p *Context) Bind(o interface{}) error {
	if err := p.Request.ParseForm(); err != nil {
		return err
	}
	if err := p.decoder.Decode(o, p.Request.Form); err != nil {
		return err
	}
	return p.validate.Struct(o)
}

// Secure https?
func (p *Context) Secure() bool {
	return p.Request.TLS != nil
}

// Locale parse locale from http-request
func (p *Context) Locale() string {
	const k = "locale"
	if p.locale != "" {
		return p.locale
	}

	lang, written := p.detectLocale(k)
	tag, _, _ := p.matcher.Match(language.Make(lang))
	if lang != tag.String() {
		written = true
		lang = tag.String()
	}
	if written {
		http.SetCookie(
			p.Writer,
			&http.Cookie{
				Name:     k,
				Value:    lang,
				Secure:   p.Secure(),
				MaxAge:   math.MaxInt32,
				HttpOnly: false,
			},
		)
	}

	p.locale = lang
	p.Payload[k] = lang
	p.Payload["languages"] = p.languages
	return lang
}
func (p *Context) detectLocale(k string) (string, bool) {
	// 1. Check URL arguments.
	if lang := p.Request.URL.Query().Get(k); lang != "" {
		return lang, true
	}

	// 2. Get language information from cookies.
	if ck, er := p.Request.Cookie(k); er == nil {
		return ck.Value, false
	}

	// 3. Get language information from 'Accept-Language'.
	return p.Request.Header.Get("Accept-Language"), true
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
func (p *Context) HTML(status int, layout, name string) {
	p.render.HTML(p.Writer, status, name, p.Payload, render.HTMLOptions{Layout: layout})
}
