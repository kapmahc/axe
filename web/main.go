package web

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
)

// Main entry
func Main() error {

	app := cli.NewApp()
	app.Name = os.Args[0]
	app.Version = fmt.Sprintf("%s (%s)", Version, BuildTime)
	app.Authors = []cli.Author{
		cli.Author{
			Name:  AuthorName,
			Email: AuthorEmail,
		},
	}
	if ts, err := time.Parse(time.RFC1123Z, BuildTime); err == nil {
		app.Compiled = ts
	}

	app.Copyright = Copyright
	app.Usage = Usage
	app.EnableBashCompletion = true
	app.Commands = make([]cli.Command, 0)

	Loop(func(p Plugin) error {
		items := p.Shell()
		app.Commands = append(app.Commands, items...)
		return nil
	})

	return app.Run(os.Args)

}
