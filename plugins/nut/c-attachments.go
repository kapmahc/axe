package nut

import (
	"github.com/facebookgo/inject"
	"github.com/kapmahc/axe/web"
	"github.com/urfave/cli"
)

// AttachmentsPlugin attachments plugin
type AttachmentsPlugin struct {
}

// Init init beans
func (p *AttachmentsPlugin) Init(*inject.Graph) error {
	return nil
}

// Shell console commands
func (p *AttachmentsPlugin) Shell() []cli.Command {
	return []cli.Command{}
}

// Mount register
func (p *AttachmentsPlugin) Mount() error {
	return nil
}

func init() {
	web.Register(&AttachmentsPlugin{})
}
