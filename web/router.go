package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter new router
func NewRouter() *Router {
	return &Router{node: mux.NewRouter()}
}

// Router http router
type Router struct {
	node *mux.Router
}

// Handler handler
func (p *Router) Handler() *mux.Router {
	return p.node
}

// Walk walk routes
func (p *Router) Walk(f mux.WalkFunc) error {
	return p.node.Walk(f)
}

// Static static assets
func (p *Router) Static(pre, dir string) {
	p.node.PathPrefix(pre).
		Handler(http.StripPrefix(pre, http.FileServer(http.Dir(dir)))).
		Methods(http.MethodGet)
}

// Group group
func (p *Router) Group(pat string) *Router {
	return &Router{node: p.node.PathPrefix(pat).Subrouter()}
}

// GET http get
func (p *Router) GET(pat string, hnd http.HandlerFunc) {
	p.node.HandleFunc(pat, hnd).Methods(http.MethodGet)
}

// POST http post
func (p *Router) POST(pat string, hnd http.HandlerFunc) {
	p.node.HandleFunc(pat, hnd).Methods(http.MethodPost)
}

// DELETE http delete
func (p *Router) DELETE(pat string, hnd http.HandlerFunc) {
	p.node.HandleFunc(pat, hnd).Methods(http.MethodDelete)
}
