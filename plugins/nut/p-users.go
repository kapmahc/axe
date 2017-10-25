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
	Jwt      *web.Jwt      `inject:""`
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

func init() {
	web.Register(&UsersPlugin{})
}
