package nut

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/kapmahc/axe/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"
)

type fmUsersSignIn struct {
}

func postUsersSignIn(l string, c *web.Context) (interface{}, error) {
	return web.H{}, nil
}

type fmUsersSignUp struct {
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"email"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"passwordConfirmation" validate:"eqfield=Password"`
}

func postUsersSignUp(l string, c *web.Context) (interface{}, error) {
	var fm fmUsersSignUp
	if err := c.Bind(&fm); err != nil {
		return nil, err
	}

	ip := c.ClientIP()
	i18n := I18N()
	if err := Tx(func(tx *pg.Tx) error {
		user, err := AddEmailUser(tx, fm.Name, fm.Email, fm.Password)
		if err != nil {
			return err
		}
		if err = AddLog(tx, user.ID, ip, i18n.T(l, "nut.logs.sign-up")); err != nil {
			return err
		}
		if err = sendEmail(c, user, actConfirm); err != nil {
			log.Error(err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return web.H{}, nil
}

const (
	actConfirm       = "confirm"
	actUnlock        = "unlock"
	actResetPassword = "reset-password"

	// SendEmailJob send email
	SendEmailJob = "send.email"
)

func sendEmail(ctx *web.Context, user *User, act string) error {
	lang := ctx.Get(web.LOCALE).(string)
	cm := jws.Claims{}
	cm.Set("act", act)
	cm.Set("uid", user.UID)
	tkn, err := JWT().Sum(cm, time.Hour*6)
	if err != nil {
		return err
	}

	obj := struct {
		Home  string
		Token string
	}{
		Home:  ctx.Home(),
		Token: string(tkn),
	}

	i18n := I18N()
	subject, err := i18n.H(lang, fmt.Sprintf("nut.emails.%s.subject", act), obj)
	if err != nil {
		return err
	}
	body, err := i18n.H(lang, fmt.Sprintf("nut.emails.%s.body", act), obj)
	if err != nil {
		return err
	}

	return JOBBER().Send(SendEmailJob, 0, map[string]string{
		"to":      user.Email,
		"subject": subject,
		"body":    body,
	})

}

func init() {
	Mount(func(rt *mux.Router) {
		unt := rt.PathPrefix("/api/users").Subrouter()
		unt.HandleFunc("/sign-in", JSON(postUsersSignIn)).Methods(http.MethodPost)
		unt.HandleFunc("/sign-up", JSON(postUsersSignUp)).Methods(http.MethodPost)
	})

	JOBBER().Register(SendEmailJob, func(id string, payload []byte) error {
		var buf bytes.Buffer
		dec := gob.NewDecoder(&buf)
		buf.Write(payload)
		arg := make(map[string]interface{})
		if err := dec.Decode(arg); err != nil {
			return err
		}

		to := arg["to"].(string)
		subject := arg["subject"].(string)
		body := arg["body"].(string)
		if viper.GetString("env") != web.PRODUCTION {
			log.Debugf("send to %s: %s\n%s", to, subject, body)
			return nil
		}

		smtp := make(map[string]interface{})
		if err := SETTINGS().Get("site.smtp", &smtp); err != nil {
			return err
		}

		sender := smtp["username"].(string)
		msg := gomail.NewMessage()
		msg.SetHeader("From", sender)
		msg.SetHeader("To", to)
		msg.SetHeader("Subject", subject)
		msg.SetBody("text/html", body)

		dia := gomail.NewDialer(
			smtp["host"].(string),
			smtp["port"].(int),
			sender,
			smtp["password"].(string),
		)

		return dia.DialAndSend(msg)

	})
}
