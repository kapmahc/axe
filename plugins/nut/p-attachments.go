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
	UEditor *web.UEditor `inject:""`
	I18n    *web.I18n    `inject:""`
	Cache   *web.Cache   `inject:""`
	Router  *gin.Engine  `inject:""`
	DB      *pg.DB       `inject:""`
	Layout  *Layout      `inject:""`
	S3      *web.S3      `inject:""`
	Jwt     *web.Jwt     `inject:""`
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
	ueditor := p.UEditor.Upload(
		p.upload, p.list(func(a *Attachment) bool {
			return a.IsPicture()
		}),
		p.list(func(a *Attachment) bool {
			return !a.IsPicture()
		}),
	)

	htm := p.Router.Group("/attachments")
	htm.GET("/ueditor", ueditor)
	htm.POST("/ueditor", ueditor)

	api := p.Router.Group("/api", p.Layout.MustSignInMiddleware)
	api.GET("/attachments", p.Layout.JSON(p.index))
	api.DELETE("/attachments/:id", p.canEdit, p.Layout.JSON(p.destroy))
	return nil
}

func init() {
	web.Register(&AttachmentsPlugin{})
}
