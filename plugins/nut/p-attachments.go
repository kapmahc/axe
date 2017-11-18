package nut

import (
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
	"github.com/urfave/cli"
)

// AttachmentsPlugin attachments plugin
type AttachmentsPlugin struct {
	I18n   *web.I18n   `inject:""`
	Cache  *web.Cache  `inject:""`
	Router *gin.Engine `inject:""`
	DB     *pg.DB      `inject:""`
	Layout *Layout     `inject:""`
	S3     *web.S3     `inject:""`
	Jwt    *web.Jwt    `inject:""`
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
	rt := p.Router
	rt.POST("/attachments", p.Layout.JSON(p.create))
	rt.GET("/attachments", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.index))
	rt.DELETE("/attachments/:id", p.Layout.MustSignInMiddleware, p.canEdit, p.Layout.JSON(p.destroy))
	return nil
}

func init() {
	web.Register(&AttachmentsPlugin{})
}
