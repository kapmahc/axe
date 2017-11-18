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
	rtn := p.Router.Group("/users")
	rtn.GET("/confirm/:token", p.Layout.Redirect("/", p.getConfirmToken))
	rtn.GET("/unlock/:token", p.Layout.Redirect("/", p.getUnlockToken))
	rtn.POST("/sign-in", p.Layout.JSON(p.postSignIn))
	rtn.POST("/sign-up", p.Layout.JSON(p.postSignUp))
	rtn.POST("/confirm", p.Layout.JSON(p.postConfirm))
	rtn.POST("/unlock", p.Layout.JSON(p.postUnlock))
	rtn.POST("/forgot-password", p.Layout.JSON(p.postForgotPassword))
	rtn.POST("/reset-password", p.Layout.JSON(p.postResetPassword))

	rtm := p.Router.Group("/users", p.Layout.MustSignInMiddleware)
	rtm.GET("/logs", p.Layout.JSON(p.getLogs))
	rtm.GET("/profile", p.Layout.JSON(p.getProfile))
	rtm.POST("/profile", p.Layout.JSON(p.postProfile))
	rtm.POST("/change-password", p.Layout.JSON(p.postChangePassword))
	rtm.DELETE("/sign-out", p.Layout.JSON(p.deleteSignOut))

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
