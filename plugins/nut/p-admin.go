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
	return nil
}

func init() {
	web.Register(&AdminPlugin{})
}
