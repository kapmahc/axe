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

	p.Layout.UEditor("/cards/summary/edit", p.checkCardToken, p.editCardH, p.updateCardH)

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

	api.GET("/links", p.Layout.JSON(p.indexLinks))
	api.GET("/links/:id", p.Layout.JSON(p.showLink))
	api.POST("/links", p.Layout.JSON(p.createLink))
	api.POST("/links/:id", p.Layout.JSON(p.updateLink))
	api.DELETE("/links/:id", p.Layout.JSON(p.destroyLink))

	api.GET("/cards", p.Layout.JSON(p.indexCards))
	api.GET("/cards/:id", p.Layout.JSON(p.showCard))
	api.POST("/cards", p.Layout.JSON(p.createCard))
	api.POST("/cards/:id", p.Layout.JSON(p.updateCard))
	api.DELETE("/cards/:id", p.Layout.JSON(p.destroyCard))

	api.GET("/friend-links", p.Layout.JSON(p.indexFriendLinks))
	api.GET("/friend-links/:id", p.Layout.JSON(p.showFriendLink))
	api.POST("/friend-links", p.Layout.JSON(p.createFriendLink))
	api.POST("/friend-links/:id", p.Layout.JSON(p.updateFriendLink))
	api.DELETE("/friend-links/:id", p.Layout.JSON(p.destroyFriendLink))

	return nil
}

func init() {
	web.Register(&AdminPlugin{})
}
