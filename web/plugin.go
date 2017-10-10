package web

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/urfave/cli"
	"github.com/urfave/negroni"
)

// Plugin plugin
type Plugin interface {
	Open(*sql.DB)
	Shell() []cli.Command
	Mount(*mux.Router, *negroni.Negroni)
}

var plugins []Plugin

// Register register plugins
func Register(args ...Plugin) {
	plugins = append(plugins, args...)
}

// Loop loop plugins
func Loop(f func(Plugin) error) error {
	for _, p := range plugins {
		if e := f(p); e != nil {
			return e
		}
	}
	return nil
}
