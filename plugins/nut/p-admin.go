package nut

import (
	"github.com/facebookgo/inject"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/web"
	"github.com/urfave/cli"
)

// AdminPlugin admin plugin
type AdminPlugin struct {
	I18n     *web.I18n     `inject:""`
	Cache    *web.Cache    `inject:""`
	Jobber   *web.Jobber   `inject:""`
	Router   *gin.Engine   `inject:""`
	Settings *web.Settings `inject:""`
	DB       *pg.DB        `inject:""`
	Redis    *redis.Pool   `inject:""`
	Dao      *Dao          `inject:""`
	Layout   *Layout       `inject:""`
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
	api := p.Router.Group("/api/admin", p.Layout.MustSignInMiddleware, p.Layout.MustAdminMiddleware)
	api.GET("/site/status", p.Layout.JSON(p.getSiteStatus))
	api.POST("/site/info", p.Layout.JSON(p.postSiteInfo))
	api.POST("/site/author", p.Layout.JSON(p.postSiteAuthor))
	api.GET("/site/seo", p.Layout.JSON(p.getSiteSeo))
	api.POST("/site/seo", p.Layout.JSON(p.postSiteSeo))
	api.GET("/site/smtp", p.Layout.JSON(p.getSiteSMTP))
	api.POST("/site/smtp", p.Layout.JSON(p.postSiteSMTP))

	api.GET("/users", p.Layout.JSON(p.indexUsers))
	api.GET("/leave-words", p.Layout.JSON(p.indexLeaveWords))
	api.DELETE("/leave-words/:id", p.Layout.JSON(p.destroyLeaveWord))

	api.GET("/locales", p.Layout.JSON(p.indexLocales))
	api.GET("/locales/:code", p.Layout.JSON(p.showLocale))
	api.POST("/locales", p.Layout.JSON(p.createLocale))
	api.DELETE("/locales/:id", p.Layout.JSON(p.destroyLocale))
	return nil
}

func init() {
	web.Register(&AdminPlugin{})
}
