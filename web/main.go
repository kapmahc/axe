package web

import (
	"fmt"
	"sort"
	"time"

	"github.com/urfave/cli"
)

var _commands []cli.Command

// RegisterCommand register conole command
func RegisterCommand(args ...cli.Command) {
	_commands = append(_commands, args...)
}

// Main entry
func Main(args ...string) error {

	app := cli.NewApp()
	app.Name = args[0]
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
	app.Commands = _commands

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app.Run(args)

}
