package nut

import (
	"bytes"
	"encoding/gob"

	"github.com/kapmahc/axe/web"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	gomail "gopkg.in/gomail.v2"
)

// Mount register
func (p *UsersPlugin) Mount() error {
	api := p.Router.Group("/api/users")
	api.POST("/sign-in", p.Layout.JSON(p.postUsersSignIn))
	api.POST("/sign-up", p.Layout.JSON(p.postUsersSignUp))
	api.POST("/confirm", p.Layout.JSON(p.postUsersConfirm))
	api.POST("/unlock", p.Layout.JSON(p.postUsersUnlock))
	api.POST("/forgot-password", p.Layout.JSON(p.postUsersForgotPassword))

	p.Jobber.Register(SendEmailJob, func(id string, payload []byte) error {
		var buf bytes.Buffer
		dec := gob.NewDecoder(&buf)
		buf.Write(payload)
		arg := make(map[string]string)
		if err := dec.Decode(&arg); err != nil {
			return err
		}

		to := arg["to"]
		subject := arg["subject"]
		body := arg["body"]
		if viper.GetString("env") != web.PRODUCTION {
			log.Debugf("send to %s: %s\n%s", to, subject, body)
			return nil
		}

		smtp := make(map[string]interface{})
		if err := p.Settings.Get("site.smtp", &smtp); err != nil {
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
	return nil
}