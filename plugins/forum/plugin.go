package forum

import (
	"github.com/facebookgo/inject"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/plugins/nut"
	"github.com/kapmahc/axe/web"
)

// Plugin plugin
type Plugin struct {
	I18n     *web.I18n     `inject:""`
	Cache    *web.Cache    `inject:""`
	Router   *gin.Engine   `inject:""`
	Settings *web.Settings `inject:""`
	Security *web.Security `inject:""`
	Jwt      *web.Jwt      `inject:""`
	DB       *pg.DB        `inject:""`
	Dao      *nut.Dao      `inject:""`
	Layout   *nut.Layout   `inject:""`
}

// Init init beans
func (p *Plugin) Init(*inject.Graph) error {
	return nil
}

// Mount register
func (p *Plugin) Mount() error {
	rt := p.Router.Group("/forum", p.Layout.MustSignInMiddleware)
	rt.GET("/articles", p.Layout.JSON(p.indexArticles))
	rt.GET("/articles/:id", p.Layout.JSON(p.showArticle))
	rt.POST("/articles", p.Layout.JSON(p.createArticle))
	rt.POST("/articles/:id", p.canEditArticle, p.Layout.JSON(p.updateArticle))
	rt.DELETE("/articles/:id", p.canEditArticle, p.Layout.JSON(p.destroyArticle))
	rt.GET("/tags", p.Layout.JSON(p.indexTags))
	rt.GET("/tags/:id", p.Layout.JSON(p.showTag))
	rt.POST("/tags", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.createTag))
	rt.POST("/tags/:id", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.updateTag))
	rt.DELETE("/tags/:id", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.destroyTag))
	rt.GET("/comments", p.Layout.JSON(p.indexComments))
	rt.GET("/comments/:id", p.Layout.JSON(p.showComment))
	rt.POST("/comments", p.Layout.JSON(p.createComment))
	rt.POST("/comments/:id", p.canEditComment, p.Layout.JSON(p.updateComment))
	rt.DELETE("/comments/:id", p.canEditComment, p.Layout.JSON(p.destroyComment))

	return nil
}

func init() {
	web.Register(&Plugin{})
}
