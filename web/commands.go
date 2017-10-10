package web

import (
	"github.com/urfave/cli"
)

var _commands []cli.Command

// Register register plugins
func Register(args ...cli.Command) {
	_commands = append(_commands, args...)
}
