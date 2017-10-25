package nut

import (
	"fmt"

	"github.com/facebookgo/inject"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
	"github.com/urfave/cli"
	"golang.org/x/text/language"
)

// UsersPlugin admin plugin
type UsersPlugin struct {
	I18n     *web.I18n     `inject:""`
	Cache    *web.Cache    `inject:""`
	Router   *web.Router   `inject:""`
	Jobber   *web.Jobber   `inject:""`
	Wrapper  *web.Wrapper  `inject:""`
	Settings *web.Settings `inject:""`
	DB       *pg.DB        `inject:""`
	Dao      *Dao          `inject:""`
	Layout   *Layout       `inject:""`
}

// Init init beans
func (p *UsersPlugin) Init(*inject.Graph) error {
	return nil
}

// Shell console commands
func (p *UsersPlugin) Shell() []cli.Command {
	return []cli.Command{
		{
			Name:    "users",
			Aliases: []string{"us"},
			Usage:   "users operations",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Aliases: []string{"l"},
					Usage:   "list users",
					Action: web.InjectAction(func(*cli.Context) error {
						var users []User
						if err := p.DB.Model(&users).Column("uid", "name", "email").Order("name ASC").Select(); err != nil {
							return err
						}
						fmt.Printf("UID\t\t\t\t\tNAME<EMAIL>\n")
						for _, u := range users {
							fmt.Printf("%s\t%s<%s>\n", u.UID, u.Name, u.Email)
						}
						return nil
					}),
				},
				{
					Name:    "role",
					Aliases: []string{"r"},
					Usage:   "apply/deny role to user",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Value: "",
							Usage: "role's name",
						},
						cli.StringFlag{
							Name:  "user, u",
							Value: "",
							Usage: "user's uid",
						},
						cli.IntFlag{
							Name:  "years, y",
							Value: 10,
							Usage: "years",
						},
						cli.BoolFlag{
							Name:  "deny, d",
							Usage: "deny mode",
						},
					},
					Action: web.InjectAction(func(c *cli.Context) error {
						uid := c.String("user")
						name := c.String("name")
						deny := c.Bool("deny")
						years := c.Int("years")

						ip := "127.0.0.1"
						lang := language.AmericanEnglish.String()

						if uid == "" || name == "" {
							cli.ShowSubcommandHelp(c)
							return nil
						}
						return web.Tx(p.DB, func(tx *pg.Tx) error {
							var user User
							if err := tx.Model(&user).
								Column("id").
								Where("uid = ?", uid).
								Limit(1).Select(); err != nil {
								return err
							}
							if deny {
								if err := p.Dao.Deny(tx, user.ID, name, DefaultResourceType, DefaultResourceID); err != nil {
									return err
								}
								if err := p.Dao.AddLog(tx, user.ID, ip, lang, "nut.logs.deny-role", name, DefaultResourceType, DefaultResourceID); err != nil {
									return err
								}
							} else {
								if err := p.Dao.Allow(tx, user.ID, name, DefaultResourceType, DefaultResourceID, years, 0, 0); err != nil {
									return err
								}
								if err := p.Dao.AddLog(tx, user.ID, ip, lang, "nut.logs.allow-role", name, DefaultResourceType, DefaultResourceID); err != nil {
									return err
								}
							}
							return nil
						})
					}),
				},
			},
		},
	}
}

// Mount register
func (p *UsersPlugin) Mount() error {
	return nil
}

func init() {
	web.Register(&UsersPlugin{})
}

// type fmUsersSignIn struct {
// }
//
// func postUsersSignIn(l string, c *web.Context) (interface{}, error) {
// 	return web.H{}, nil
// }
//
// type fmUsersSignUp struct {
// 	Name                 string `json:"name" validate:"required"`
// 	Email                string `json:"email" validate:"email"`
// 	Password             string `json:"password" validate:"required"`
// 	PasswordConfirmation string `json:"passwordConfirmation" validate:"eqfield=Password"`
// }
//
// func postUsersSignUp(l string, c *web.Context) (interface{}, error) {
// 	var fm fmUsersSignUp
// 	if err := c.Bind(&fm); err != nil {
// 		return nil, err
// 	}
//
// 	ip := c.ClientIP()
// 	i18n := I18N()
// 	if err := Tx(func(tx *pg.Tx) error {
// 		user, err := AddEmailUser(tx, fm.Name, fm.Email, fm.Password)
// 		if err != nil {
// 			return err
// 		}
// 		if err = AddLog(tx, user.ID, ip, i18n.T(l, "nut.logs.sign-up")); err != nil {
// 			return err
// 		}
// 		if err = sendEmail(c, user, actConfirm); err != nil {
// 			log.Error(err)
// 		}
// 		return nil
// 	}); err != nil {
// 		return nil, err
// 	}
//
// 	return web.H{}, nil
// }
//
// const (
// 	actConfirm       = "confirm"
// 	actUnlock        = "unlock"
// 	actResetPassword = "reset-password"
//
// 	// SendEmailJob send email
// 	SendEmailJob = "send.email"
// )
//
// func sendEmail(ctx *web.Context, user *User, act string) error {
// 	lang := ctx.Get(web.LOCALE).(string)
// 	cm := jws.Claims{}
// 	cm.Set("act", act)
// 	cm.Set("uid", user.UID)
// 	tkn, err := JWT().Sum(cm, time.Hour*6)
// 	if err != nil {
// 		return err
// 	}
//
// 	obj := struct {
// 		Home  string
// 		Token string
// 	}{
// 		Home:  ctx.Home(),
// 		Token: string(tkn),
// 	}
//
// 	i18n := I18N()
// 	subject, err := i18n.H(lang, fmt.Sprintf("nut.emails.%s.subject", act), obj)
// 	if err != nil {
// 		return err
// 	}
// 	body, err := i18n.H(lang, fmt.Sprintf("nut.emails.%s.body", act), obj)
// 	if err != nil {
// 		return err
// 	}
//
// 	return JOBBER().Send(SendEmailJob, 0, map[string]string{
// 		"to":      user.Email,
// 		"subject": subject,
// 		"body":    body,
// 	})
//
// }
//
// func init() {
// 	Mount(func(rt *mux.Router) {
// 		unt := rt.PathPrefix("/api/users").Subrouter()
// 		unt.HandleFunc("/sign-in", JSON(postUsersSignIn)).Methods(http.MethodPost)
// 		unt.HandleFunc("/sign-up", JSON(postUsersSignUp)).Methods(http.MethodPost)
// 	})
//
// 	JOBBER().Register(SendEmailJob, func(id string, payload []byte) error {
// 		var buf bytes.Buffer
// 		dec := gob.NewDecoder(&buf)
// 		buf.Write(payload)
// 		arg := make(map[string]interface{})
// 		if err := dec.Decode(arg); err != nil {
// 			return err
// 		}
//
// 		to := arg["to"].(string)
// 		subject := arg["subject"].(string)
// 		body := arg["body"].(string)
// 		if viper.GetString("env") != web.PRODUCTION {
// 			log.Debugf("send to %s: %s\n%s", to, subject, body)
// 			return nil
// 		}
//
// 		smtp := make(map[string]interface{})
// 		if err := SETTINGS().Get("site.smtp", &smtp); err != nil {
// 			return err
// 		}
//
// 		sender := smtp["username"].(string)
// 		msg := gomail.NewMessage()
// 		msg.SetHeader("From", sender)
// 		msg.SetHeader("To", to)
// 		msg.SetHeader("Subject", subject)
// 		msg.SetBody("text/html", body)
//
// 		dia := gomail.NewDialer(
// 			smtp["host"].(string),
// 			smtp["port"].(int),
// 			sender,
// 			smtp["password"].(string),
// 		)
//
// 		return dia.DialAndSend(msg)
//
// 	})
// }
