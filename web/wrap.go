package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

const (
	// NOTICE notice
	NOTICE = "notice"
	// WARNING warning
	WARNING = "warning"
	// ERROR error
	ERROR = "error"

	// TITLE title
	TITLE = "title"
	// MESSAGE message
	MESSAGE = "message"

	// APPLICATION application layout
	APPLICATION = "layouts/application/index"
	// DASHBOARD dashboard layout
	DASHBOARD = "layouts/dashboard/index"
)

// HTMLHandlerFunc html handler func
type HTMLHandlerFunc func(string, *Context) (H, error)

// RedirectHandlerFunc redirect handle func
type RedirectHandlerFunc func(string, *Context) error

// ObjectHandlerFunc object handle func
type ObjectHandlerFunc func(string, *Context) (interface{}, error)

// Redirect redirect
func Redirect(to string, fn RedirectHandlerFunc) HandlerFunc {
	return func(c *Context) {
		if err := fn(c.Locale(), c); err != nil {
			log.Error(err)
			ss := c.Session()
			ss.AddFlash(err.Error(), ERROR)
			c.Save(ss)
		}
		c.Redirect(http.StatusFound, to)
	}
}

// JSON render json
func JSON(fn ObjectHandlerFunc) HandlerFunc {
	return func(c *Context) {
		if val, err := fn(c.Locale(), c); err == nil {
			c.JSON(http.StatusOK, val)
		} else {
			status, body := detectError(err)
			c.Text(status, body)
		}
	}
}

func detectError(e error) (int, string) {
	log.Error(e)
	if er, ok := e.(validator.ValidationErrors); ok {
		var ss []string
		for _, it := range er {
			ss = append(ss, fmt.Sprintf("Validation for '%s' failed on the '%s' tag;", it.Field(), it.Tag()))
		}
		return http.StatusBadRequest, strings.Join(ss, "\n")
	}
	return http.StatusInternalServerError, e.Error()
}

// XML wrap xml
func XML(fn ObjectHandlerFunc) HandlerFunc {
	return func(c *Context) {
		if val, err := fn(c.Locale(), c); err == nil {
			c.XML(http.StatusOK, val)
		} else {
			status, body := detectError(err)
			c.Text(status, body)
		}
	}
}

// HTML wrap html
func HTML(layout, name string, handler HTMLHandlerFunc) HandlerFunc {
	return func(c *Context) {
		lang := c.Locale()
		flashes := H{}
		ss := c.Session()
		for _, n := range []string{NOTICE, WARNING, ERROR} {
			flashes[n] = ss.Flashes(n)
		}
		c.Save(ss)

		payload, err := handler(lang, c)
		if err != nil {
			payload = H{}
		}
		payload["locale"] = lang
		payload["flashes"] = flashes
		payload["session"] = c.Session().Values
		payload[csrf.TemplateTag] = csrf.TemplateField(c.Request)
		payload["csrfToken"] = csrf.Token(c.Request)
		if err == nil {
			c.HTML(http.StatusOK, layout, name, payload)
		} else {
			status, body := detectError(err)
			payload["reason"] = body
			payload["createdAt"] = time.Now()
			c.HTML(status, layout, "error", payload)
		}
	}
}
