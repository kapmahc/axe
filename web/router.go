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
func NewRouter(secure bool, secret []byte, theme string, helpers template.FuncMap, matcher language.Matcher) *Router {
	store := sessions.NewCookieStore(secret)

	router := mux.NewRouter()
	for k, v := range map[string]string{
		"3rd":    "node_modules",
		"assets": path.Join("themes", theme, "assets"),
	} {
		pat := "/" + k + "/"
		router.PathPrefix(pat).
			Handler(http.StripPrefix(pat, http.FileServer(http.Dir(v)))).
			Methods(http.MethodGet)
	}

	return &Router{
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
		matcher:  matcher,
		decoder:  form.NewDecoder(),
		validate: validator.New(),
		store:    store,
	}
}

// Router http router
type Router struct {
	router   *mux.Router
	csrf     func(http.Handler) http.Handler
	render   *render.Render
	matcher  language.Matcher
	validate *validator.Validate
	decoder  *form.Decoder
	store    sessions.Store
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
func (p *Router) Walk(f func(string, ...string) error) error {
	return p.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pattern, err := route.GetPathTemplate()
		// '/articles/{id}' will be '^/articles/(?P<v0>[^/]+)$'
		// name, err := route.GetPathRegexp() {"^surname=(?P<v0>.*)$}
		// 'Queries("surname", "{surname}") will return
		// name, err := route.GetQueriesTemplates()
		// name, err := route.GetQueriesRegexp()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return err
		}
		if len(methods) == 0 {
			return nil
		}
		return f(pattern, methods...)
	})
}

// Group group
func (p *Router) Group(prefix string) *Router {
	return &Router{
		render:   p.render,
		csrf:     p.csrf,
		validate: p.validate,
		decoder:  p.decoder,
		matcher:  p.matcher,
		router:   p.router.PathPrefix(prefix).Subrouter(),
		store:    p.store,
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

// Form handle form
func (p *Router) Form(
	pattern,
	layout, name string, get HTMLHandlerFunc,
	post ObjectHandlerFunc) {
	p.GET(pattern, HTML(layout, name, get))
	p.POST(pattern, JSON(post))
}

// Crud handle crud
func (p *Router) Crud(
	prefix,
	layout, name string,
	index HTMLHandlerFunc,
	_new HTMLHandlerFunc, create ObjectHandlerFunc,
	show HTMLHandlerFunc,
	edit HTMLHandlerFunc, update ObjectHandlerFunc,
	destroy ObjectHandlerFunc) {
	if index != nil {
		p.GET(prefix, HTML(layout, name+"/index", index))
	}
	if _new != nil {
		p.GET(prefix+"/new", HTML(layout, name+"/new", _new))
	}
	if create != nil {
		p.POST(prefix+"", JSON(create))
	}
	if show != nil {
		p.GET(prefix+"/{id}", HTML(layout, name+"/show", show))
	}
	if edit != nil {
		p.GET(prefix+"/{id}/edit", HTML(layout, name+"/edit", edit))
	}
	if update != nil {
		p.POST(prefix+"/{id}", JSON(update))
	}
	if destroy != nil {
		p.DELETE(prefix+"/{id}", JSON(destroy))
	}
}

func (p *Router) add(method, pattern string, handler HandlerFunc) {
	p.router.
		HandleFunc(pattern, func(wrt http.ResponseWriter, req *http.Request) {
			handler(NewContext(req, wrt, p.store, p.matcher, p.render, p.validate, p.decoder))
		}).
		Methods(method)
}
