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
	htm := p.Router.Group("/users")
	htm.GET("/confirm/:token", p.Layout.Redirect("/", p.getConfirmToken))
	htm.GET("/unlock/:token", p.Layout.Redirect("/", p.getUnlockToken))

	api := p.Router.Group("/api/users")
	api.POST("/sign-in", p.Layout.JSON(p.postSignIn))
	api.POST("/sign-up", p.Layout.JSON(p.postSignUp))
	api.POST("/confirm", p.Layout.JSON(p.postConfirm))
	api.POST("/unlock", p.Layout.JSON(p.postUnlock))
	api.POST("/forgot-password", p.Layout.JSON(p.postForgotPassword))
	api.POST("/reset-password", p.Layout.JSON(p.postResetPassword))

	apiM := p.Router.Group("/api/users", p.Layout.MustSignInMiddleware)
	apiM.GET("/logs", p.Layout.JSON(p.getLogs))
	apiM.GET("/profile", p.Layout.JSON(p.getProfile))
	apiM.POST("/profile", p.Layout.JSON(p.postProfile))
	apiM.POST("/change-password", p.Layout.JSON(p.postChangePassword))
	apiM.DELETE("/sign-out", p.Layout.JSON(p.deleteSignOut))

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
