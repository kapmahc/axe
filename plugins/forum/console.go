package forum

import (
	"encoding/csv"
	"os"
	"time"

	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/plugins/nut"
	"github.com/kapmahc/axe/web"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	blackfriday "gopkg.in/russross/blackfriday.v2"
)

// Shell console commands
func (p *Plugin) Shell() []cli.Command {
	return []cli.Command{
		{
			Name:  "forum",
			Usage: "forum operations",
			Subcommands: []cli.Command{
				{
					Name:    "import",
					Aliases: []string{"i"},
					Usage:   "Import articles: copy (select title,type,body from forum_articles) to '/tmp/articles.csv' (format CSV);",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "csv file name",
						},
						cli.StringFlag{
							Name:  "user, u",
							Usage: "user's uid",
						},
					},
					Action: web.InjectAction(func(c *cli.Context) error {
						name := c.String("name")
						if name == "" {
							cli.ShowSubcommandHelp(c)
							return nil
						}
						log.Info("import articles from file ", name)
						fd, err := os.Open(name)
						if err != nil {
							return err
						}
						defer fd.Close()
						items, err := csv.NewReader(fd).ReadAll()
						if err != nil {
							return err
						}
						now := time.Now()
						return p.DB.RunInTransaction(func(tx *pg.Tx) error {
							var u nut.User
							if err := tx.Model(&u).Column("id").
								Where("uid = ?", c.String("user")).
								Limit(1).Select(); err != nil {
								return err
							}
							for _, it := range items {
								a := Article{
									Title:     it[0],
									Type:      web.TypeHTML,
									UserID:    u.ID,
									UpdatedAt: now,
								}

								if it[1] == web.TypeMARKDOWN {
									a.Body = string(blackfriday.Run([]byte(it[2])))
								} else {
									a.Body = it[2]
								}
								if err := tx.Insert(&a); err != nil {
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
