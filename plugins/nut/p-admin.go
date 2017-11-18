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

	rt := p.Router.Group("/admin", p.Layout.MustSignInMiddleware, p.Layout.MustAdminMiddleware)
	rt.GET("/site/status", p.Layout.JSON(p.getSiteStatus))
	rt.POST("/site/info", p.Layout.JSON(p.postSiteInfo))
	rt.POST("/site/author", p.Layout.JSON(p.postSiteAuthor))
	rt.GET("/site/seo", p.Layout.JSON(p.getSiteSeo))
	rt.POST("/site/seo", p.Layout.JSON(p.postSiteSeo))
	rt.GET("/site/smtp", p.Layout.JSON(p.getSiteSMTP))
	rt.POST("/site/smtp", p.Layout.JSON(p.postSiteSMTP))
	rt.GET("/site/home", p.Layout.JSON(p.getSiteHome))
	rt.POST("/site/home", p.Layout.JSON(p.postSiteHome))

	rt.GET("/users", p.Layout.JSON(p.indexUsers))
	rt.GET("/leave-words", p.Layout.JSON(p.indexLeaveWords))
	rt.DELETE("/leave-words/:id", p.Layout.JSON(p.destroyLeaveWord))

	rt.GET("/locales", p.Layout.JSON(p.indexLocales))
	rt.GET("/locales/:code", p.Layout.JSON(p.showLocale))
	rt.POST("/locales", p.Layout.JSON(p.createLocale))
	rt.DELETE("/locales/:id", p.Layout.JSON(p.destroyLocale))

	rt.GET("/links", p.Layout.JSON(p.indexLinks))
	rt.GET("/links/:id", p.Layout.JSON(p.showLink))
	rt.POST("/links", p.Layout.JSON(p.createLink))
	rt.POST("/links/:id", p.Layout.JSON(p.updateLink))
	rt.DELETE("/links/:id", p.Layout.JSON(p.destroyLink))

	rt.GET("/cards", p.Layout.JSON(p.indexCards))
	rt.GET("/cards/:id", p.Layout.JSON(p.showCard))
	rt.POST("/cards", p.Layout.JSON(p.createCard))
	rt.POST("/cards/:id", p.Layout.JSON(p.updateCard))
	rt.DELETE("/cards/:id", p.Layout.JSON(p.destroyCard))

	rt.GET("/friend-links", p.Layout.JSON(p.indexFriendLinks))
	rt.GET("/friend-links/:id", p.Layout.JSON(p.showFriendLink))
	rt.POST("/friend-links", p.Layout.JSON(p.createFriendLink))
	rt.POST("/friend-links/:id", p.Layout.JSON(p.updateFriendLink))
	rt.DELETE("/friend-links/:id", p.Layout.JSON(p.destroyFriendLink))

	return nil
}

func init() {
	web.Register(&AdminPlugin{})
}
