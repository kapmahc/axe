package web

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/go-playground/form"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/unrolled/render"
	"golang.org/x/text/language"
	validator "gopkg.in/go-playground/validator.v9"
)

// NewRouter new router
func NewRouter(secure bool, secret []byte, theme string, helpers template.FuncMap, langs ...language.Tag) *Router {
	store := sessions.NewCookieStore(secret)

	router := mux.NewRouter()
	for k, v := range map[string]string{
		"3rd":    "node_modules",
		"assets": path.Join("themes", theme, "assets"),
	} {
		pat := "/" + k + "/"
		router.PathPrefix(pat).
			Handler(http.StripPrefix(pat, http.FileServer(http.Dir(v)))).
			Methods(http.MethodGet).
			Name(k)
	}

	return &Router{
		name: "",
		render: render.New(render.Options{
			Directory:  path.Join("themes", theme, "views"),
			Extensions: []string{".html"},
			Funcs:      []template.FuncMap{helpers},
		}),

		router: router,
		csrf: csrf.Protect(
			secret,
			csrf.Secure(secure),
			csrf.CookieName("csrf"),
			csrf.RequestHeader("Authenticity-Token"),
			csrf.FieldName("authenticity_token"),
		),
		matcher:   language.NewMatcher(langs),
		decoder:   form.NewDecoder(),
		validate:  validator.New(),
		store:     store,
		languages: langs,
	}
}

// Router http router
type Router struct {
	name      string
	router    *mux.Router
	csrf      func(http.Handler) http.Handler
	render    *render.Render
	matcher   language.Matcher
	validate  *validator.Validate
	decoder   *form.Decoder
	languages []language.Tag
	store     sessions.Store
}

// Start start server
func (p *Router) Start(port int, grace bool) error {
	log.Infof(
		"application starting on http://localhost:%d",
		port,
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: p.Handler(),
	}

	if !grace {
		return srv.ListenAndServe()
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Warn("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	log.Warn("server exiting")
	return nil
}

// Handler http handler
func (p *Router) Handler() http.Handler {
	return p.csrf(p.router)
}

// Walk walk routes
func (p *Router) Walk(f func(string, string, ...string) error) error {
	return p.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		// '/articles/{id}' will be '^/articles/(?P<v0>[^/]+)$'
		// tpl, err := route.GetPathRegexp() {"^surname=(?P<v0>.*)$}
		// 'Queries("surname", "{surname}") will return
		// tpl, err := route.GetQueriesTemplates()
		// tpl, err := route.GetQueriesRegexp()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return err
		}
		return f(route.GetName(), tpl, methods...)
	})
}

// Group group
func (p *Router) Group(name, prefix string) *Router {
	if p.name != "" {
		name = p.name + "." + name
	}
	return &Router{
		name:      name,
		render:    p.render,
		csrf:      p.csrf,
		validate:  p.validate,
		decoder:   p.decoder,
		matcher:   p.matcher,
		router:    p.router.PathPrefix(prefix).Subrouter(),
		languages: p.languages,
		store:     p.store,
	}
}

// GET http get
func (p *Router) GET(pattern string, handler HandlerFunc) {
	p.add(http.MethodGet, pattern, handler)
}

// POST http post
func (p *Router) POST(pattern string, handler HandlerFunc) {
	p.add(http.MethodPost, pattern, handler)
}

// PATCH http patch
func (p *Router) PATCH(pattern string, handler HandlerFunc) {
	p.add(http.MethodPatch, pattern, handler)
}

// DELETE http delete
func (p *Router) DELETE(pattern string, handler HandlerFunc) {
	p.add(http.MethodDelete, pattern, handler)
}

func (p *Router) add(method, pattern string, handler HandlerFunc) {
	p.router.
		HandleFunc(pattern, func(wrt http.ResponseWriter, req *http.Request) {
			handler(NewContext(req, wrt, p.store, p.matcher, p.render, p.validate, p.decoder))
		}).
		Methods(method).
		Name(p.name)
}
