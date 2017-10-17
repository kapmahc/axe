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

var (
	_router = mux.NewRouter()
)

// Mount web mount points
func Mount(f func(*mux.Router)) {
	f(_router)
}

func listen() error {
	port := viper.GetInt("server.port")
	addr := fmt.Sprintf(":%d", port)
	log.Infof(
		"application starting on http://localhost:%d",
		port,
	)
	hnd, err := httpServer()
	if err != nil {
		return err
	}
	srv := &http.Server{
		Addr: addr,
		Handler: csrf.Protect(
			[]byte(viper.GetString("secret")),
			csrf.CookieName("csrf"),
			csrf.RequestHeader("Authenticity-Token"),
			csrf.FieldName("authenticity_token"),
			csrf.Secure(viper.GetBool("server.ssl")),
		)(hnd),
	}

	if viper.GetString("env") != web.PRODUCTION {
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

func httpServer() (http.Handler, error) {
	for k, v := range map[string]string{
		"3rd":    "node_modules",
		"assets": path.Join("themes", viper.GetString("server.theme"), "assets"),
	} {
		pre := "/" + k + "/"
		_router.PathPrefix(pre).
			Handler(http.StripPrefix(pre, http.FileServer(http.Dir(v)))).
			Methods(http.MethodGet)
	}

	ng := negroni.New()
	ng.Use(negroni.NewRecovery())
	ng.UseFunc(web.LoggerMiddleware)
	i18n, err := I18N().Middleware()
	if err != nil {
		return nil, err
	}
	ng.UseFunc(i18n)
	ng.UseHandler(_router)
	return ng, nil
}
