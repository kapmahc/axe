package pos

import (
	"github.com/facebookgo/inject"
	"github.com/kapmahc/axe/web"
	"github.com/urfave/cli"
)

// Plugin plugin
type Plugin struct {
}

// Init init beans
func (p *Plugin) Init(*inject.Graph) error {
	return nil
}

// Shell console commands
func (p *Plugin) Shell() []cli.Command {
	return []cli.Command{}
}

// Mount register
func (p *Plugin) Mount() error {
	return nil
}

func init() {
	web.Register(&Plugin{})
}
