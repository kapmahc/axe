package nut

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/kapmahc/axe/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"
)

// Plugin plugin
type Plugin interface {
	Mount(*mux.Router)
}

var (
	_plugins []Plugin
)

// Register register plugin
func Register(args ...Plugin) {
	_plugins = append(_plugins, args...)
}

func listen() error {
	port := viper.GetInt("server.port")
	addr := fmt.Sprintf(":%d", port)
	log.Infof(
		"application starting on http://localhost:%d",
		port,
	)
	srv := &http.Server{
		Addr: addr,
		Handler: csrf.Protect(
			[]byte(viper.GetString("secret")),
			csrf.CookieName("csrf"),
			csrf.RequestHeader("Authenticity-Token"),
			csrf.FieldName("authenticity_token"),
			csrf.Secure(viper.GetBool("server.ssl")),
		)(httpServer()),
	}

	if viper.GetString("env") != "production" {
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

func httpServer() http.Handler {
	rt := mux.NewRouter()
	for _, p := range _plugins {
		p.Mount(rt)
	}
	for k, v := range map[string]string{
		"3rd":    "node_modules",
		"assets": path.Join("themes", viper.GetString("server.theme"), "assets"),
	} {
		pre := "/" + k + "/"
		rt.PathPrefix(pre).Handler(http.StripPrefix(pre, http.FileServer(http.Dir(v)))).Methods(http.MethodGet)
	}

	ng := negroni.New()
	ng.Use(negroni.NewRecovery())
	ng.UseFunc(web.LoggerMiddleware)
	ng.UseFunc(I18N().Middleware())
	ng.UseHandler(rt)
	return ng
}
