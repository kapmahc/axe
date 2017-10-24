package nut

import (
	"github.com/facebookgo/inject"
	"github.com/kapmahc/axe/web"
	"github.com/urfave/cli"
)

// AdminPlugin admin plugin
type AdminPlugin struct {
}

// Init init beans
func (p *AdminPlugin) Init(*inject.Graph) error {
	return nil
}

// Shell console commands
func (p *AdminPlugin) Shell() []cli.Command {
	return []cli.Command{}
}

// Mount register
func (p *AdminPlugin) Mount() error {
	return nil
}

func init() {
	web.Register(&AdminPlugin{})
}
