package nut

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/facebookgo/inject"
	"github.com/gorilla/mux"
	"github.com/kapmahc/axe/web"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// HomePlugin admin plugin
type HomePlugin struct {
}

// Init init beans
func (p *HomePlugin) Init(*inject.Graph) error {
	return nil
}

// Shell console commands
func (p *HomePlugin) Shell() []cli.Command {
	return []cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate file template",
			Subcommands: []cli.Command{
				{
					Name:    "config",
					Aliases: []string{"c"},
					Usage:   "generate config file",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "environment, e",
							Value: "development",
							Usage: "environment, like: development, production, stage, test...",
						},
					},
					Action: generateConfig,
				},
				{
					Name:    "nginx",
					Aliases: []string{"ng"},
					Usage:   "generate nginx.conf",
					Action:  Open(generateNginxConf, false),
				},
				{
					Name:    "openssl",
					Aliases: []string{"ssl"},
					Usage:   "generate ssl certificates",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "name",
						},
						cli.StringFlag{
							Name:  "country, c",
							Value: "Earth",
							Usage: "country",
						},
						cli.StringFlag{
							Name:  "organization, o",
							Value: "Mother Nature",
							Usage: "organization",
						},
						cli.IntFlag{
							Name:  "years, y",
							Value: 1,
							Usage: "years",
						},
					},
					Action: generateSsl,
				},
				{
					Name:    "migration",
					Usage:   "generate migration file",
					Aliases: []string{"m"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "name",
						},
					},
					Action: generateMigration,
				},
				{
					Name:    "locale",
					Usage:   "generate locale file",
					Aliases: []string{"l"},
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "locale name",
						},
					},
					Action: generateLocale,
				},
			},
		},
		{
			Name:    "cache",
			Aliases: []string{"c"},
			Usage:   "cache operations",
			Subcommands: []cli.Command{
				{
					Name:    "list",
					Usage:   "list all cache keys",
					Aliases: []string{"l"},
					Action: Open(func(_ *cli.Context) error {
						keys, err := CACHE().Keys()
						if err != nil {
							return err
						}
						for _, k := range keys {
							fmt.Println(k)
						}
						return nil
					}, true),
				},
				{
					Name:    "clear",
					Usage:   "clear cache items",
					Aliases: []string{"c"},
					Action: Open(func(_ *cli.Context) error {
						return CACHE().Flush()
					}, true),
				},
			},
		},
		{
			Name:    "database",
			Aliases: []string{"db"},
			Usage:   "database operations",
			Subcommands: []cli.Command{
				{
					Name:    "example",
					Usage:   "scripts example for create database and user",
					Aliases: []string{"e"},
					Action:  Open(databaseExample, false),
				},
				{
					Name:    "migrate",
					Usage:   "migrate the DB to the most recent version available",
					Aliases: []string{"m"},
					Action:  Open(runDatabase("up"), false),
				},
				{
					Name:    "rollback",
					Usage:   "roll back the version by 1",
					Aliases: []string{"r"},
					Action:  Open(runDatabase("down"), false),
				},
				{
					Name:    "version",
					Usage:   "dump the migration status for the current DB",
					Aliases: []string{"v"},
					Action:  Open(runDatabase("version"), false),
				},
				{
					Name:    "connect",
					Usage:   "connect database",
					Aliases: []string{"c"},
					Action:  Open(connectDatabase, false),
				},
				{
					Name:    "create",
					Usage:   "create database",
					Aliases: []string{"n"},
					Action:  Open(createDatabase, false),
				},
				{
					Name:    "drop",
					Usage:   "drop database",
					Aliases: []string{"d"},
					Action:  Open(dropDatabase, false),
				},
			},
		},
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "start the app server",
			Action: Open(func(_ *cli.Context) error {
				go func() {
					// ----------
					host, err := os.Hostname()
					if err != nil {
						log.Error(err)
					}
					for {
						if err := JOBBER().Receive(host); err != nil {
							log.Error(err)
							time.Sleep(5 * time.Second)
						}
					}
				}()
				// -------
				return listen()
			}, true),
		},
		{
			Name:    "routes",
			Aliases: []string{"rt"},
			Usage:   "print out all defined routes",
			Action: func(_ *cli.Context) error {
				tpl := "%-7s %s\n"
				fmt.Printf(tpl, "METHOD", "PATH")
				return _router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
					pat, err := route.GetPathTemplate()
					if err != nil {
						return err
					}
					mtd, err := route.GetMethods()
					if err != nil {
						return err
					}
					if len(mtd) == 0 {
						return nil
					}
					fmt.Printf(tpl, strings.Join(mtd, ","), pat)
					return nil
				})
			},
		},
	}
}

// Mount register
func (p *HomePlugin) Mount() error {
	return nil
}

func init() {
	web.Register(&HomePlugin{})
}

// func getAPISiteInfo(l string, c *web.Context) (interface{}, error) {
// 	i18n := I18N()
// 	// -----------
// 	langs, err := i18n.Languages()
// 	if err != nil {
// 		return nil, err
// 	}
// 	data := web.H{"locale": l, "languages": langs}
// 	// -----------
// 	for _, k := range []string{"title", "subhead", "keywords", "description", "copyright"} {
// 		data[k] = i18n.T(l, "site."+k)
// 	}
// 	// -----------
// 	author := web.H{}
// 	for _, k := range []string{"name", "email"} {
// 		author[k] = i18n.T(l, "site.author."+k)
// 	}
// 	data["author"] = author
// 	return data, nil
// }
//
// type fmInstall struct {
// 	Title                string `json:"title" validate:"required"`
// 	Subhead              string `json:"subhead" validate:"required"`
// 	Name                 string `json:"name" validate:"required"`
// 	Email                string `json:"email" validate:"email"`
// 	Password             string `json:"password" validate:"required"`
// 	PasswordConfirmation string `json:"passwordConfirmation" validate:"eqfield=Password"`
// }
//
// func postInstall(l string, c *web.Context) (interface{}, error) {
// 	var fm fmInstall
// 	if err := c.Bind(&fm); err != nil {
// 		return nil, err
// 	}
//
// 	now := time.Now()
// 	ip := c.ClientIP()
// 	i18n := I18N()
// 	if err := Tx(func(tx *pg.Tx) error {
// 		cnt, err := tx.Model(&User{}).Count()
// 		if err != nil {
// 			return err
// 		}
// 		if cnt > 0 {
// 			return i18n.E(l, "errors.forbidden")
// 		}
// 		if err = i18n.Set(tx, l, "site.title", fm.Title); err != nil {
// 			return err
// 		}
// 		if err = i18n.Set(tx, l, "site.subhead", fm.Subhead); err != nil {
// 			return err
// 		}
// 		user, err := AddEmailUser(tx, fm.Name, fm.Email, fm.Password)
// 		if err != nil {
// 			return err
// 		}
// 		if err = AddLog(tx, user.ID, ip, i18n.T(l, "nut.logs.sign-up")); err != nil {
// 			return err
// 		}
// 		user.ConfirmedAt = &now
// 		user.UpdatedAt = now
// 		if _, err = tx.Model(user).Column("confirmed_at", "updated_at").Update(); err != nil {
// 			return err
// 		}
// 		for _, rn := range []string{RoleRoot, RoleAdmin} {
// 			if err := Allow(tx, user.ID, rn, DefaultResourceType, DefaultResourceID, 50, 0, 0); err != nil {
// 				return err
// 			}
// 			if err := AddLog(tx, user.ID, ip, i18n.T(l, "nut.logs.apply-role", rn, DefaultResourceType, DefaultResourceID)); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	}); err != nil {
// 		return nil, err
// 	}
// 	return web.H{}, nil
// }
//
// func init() {
// 	Mount(func(rt *mux.Router) {
// 		api := rt.PathPrefix("/api").Subrouter()
// 		api.HandleFunc("/site/info", JSON(getAPISiteInfo)).Methods(http.MethodGet)
// 		api.HandleFunc("/install", JSON(postInstall)).Methods(http.MethodPost)
// 	})
// }
