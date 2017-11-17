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
	const signInURL = "/users/sign-in"
	rt := p.Router.Group("/users")
	rt.DELETE("/sign-out", web.JSON(p.deleteSignOut))
	rt.GET("/logs", web.HTML(web.APPLICATION, "nut/users/logs", p.getLogs))
	rt.Form("/profile", web.DASHBOARD, "nut/users/profile", p.getProfile, p.postProfile)
	rt.Form("/change-password", web.APPLICATION, "nut/users/change-password", p.getChangePassword, p.postChangePassword)

	rt.Form("/sign-in", web.APPLICATION, "nut/users/sign-in", p.getSignIn, p.postSignIn)
	rt.Form("/sign-up", web.APPLICATION, "nut/users/sign-up", p.getSignUp, p.postSignUp)
	rt.Form("/confirm", web.APPLICATION, "nut/users/confirm", p.getConfirm, p.postConfirm)
	rt.GET("/confirm/{token}", web.Redirect(signInURL, p.getConfirmToken))
	rt.Form("/unlock", web.APPLICATION, "nut/users/unlock", p.getUnlock, p.postUnlock)
	rt.GET("/unlock/{token}", web.Redirect(signInURL, p.getUnlockToken))
	rt.Form("/forgot-password", web.APPLICATION, "nut/users/forgot-password", p.getForgotPassword, p.postForgotPassword)
	rt.Form("/reset-password/{token}", web.APPLICATION, "nut/users/reset-password", p.getResetPassword, p.postResetPassword)

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
