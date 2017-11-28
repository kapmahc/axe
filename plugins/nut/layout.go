package nut

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/kapmahc/axe/web"
	log "github.com/sirupsen/logrus"
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
	APPLICATION = "layouts/application"
	// DASHBOARD dashboard layout
	DASHBOARD = "layouts/dashboard"

	BANNER = `
   ____     __     __    _____
  (    )   (_ \   / _)  / ___/
  / /\ \     \ \_/ /   ( (__
 ( (__) )     \   /     ) __)
  )    (      / _ \    ( (
 /  /\  \   _/ / \ \_   \ \___
/__(  )__\ (__/   \__)   \____\
`
)

// HTMLHandlerFunc html handler func
type HTMLHandlerFunc func(string, *gin.Context) (gin.H, error)

// RedirectHandlerFunc redirect handle func
type RedirectHandlerFunc func(string, *gin.Context) error

// ObjectHandlerFunc object handle func
type ObjectHandlerFunc func(string, *gin.Context) (interface{}, error)

// Layout layout
type Layout struct {
	Store sessions.Store `inject:""`
}

func (p *Layout) Session(c *gin.Context) *sessions.Session {
	ss, _ := p.Store.Get(c.Request, "session")
	return ss
}
func (p *Layout) Save(s *sessions.Session, c *gin.Context) {
	if err := s.Save(c.Request, c.Writer); err != nil {
		log.Error(err)
	}
}

// Redirect redirect
func (p *Layout) Redirect(to string, fn RedirectHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c.MustGet(web.LOCALE).(string), c); err != nil {
			log.Error(err)
			s := p.Session(c)
			s.AddFlash(err.Error(), ERROR)
			p.Save(s, c)
		}
		c.Redirect(http.StatusFound, to)
	}
}

// JSON render json
func JSON(fn ObjectHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if val, err := fn(c.MustGet(LOCALE).(string), c); err == nil {
			c.JSON(http.StatusOK, val)
		} else {
			status, body := detectError(err)
			c.String(status, body)
		}
	}
}
