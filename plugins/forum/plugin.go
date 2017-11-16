package forum

import (
	"github.com/facebookgo/inject"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/plugins/nut"
	"github.com/kapmahc/axe/web"
)

// Plugin plugin
type Plugin struct {
	I18n     *web.I18n     `inject:""`
	Cache    *web.Cache    `inject:""`
	Settings *web.Settings `inject:""`
	Security *web.Security `inject:""`
	Jwt      *web.Jwt      `inject:""`
	DB       *pg.DB        `inject:""`
	Dao      *nut.Dao      `inject:""`
}

// Init init beans
func (p *Plugin) Init(*inject.Graph) error {
	return nil
}

// Mount register
func (p *Plugin) Mount() error {
	// p.Layout.UEditor("/forum/articles/body/edit", p.checkArticleToken, p.editArticleH, p.updateArticleH)
	// p.Layout.UEditor("/forum/comments/body/edit", p.checkCommentToken, p.editCommentH, p.updateCommentH)
	//
	// htm := p.Router.Group("/forum")
	// htm.GET("/tags", p.Layout.Application("forum-tags-index", p.indexTagsH))
	// htm.GET("/tags/show/:id", p.Layout.Application("forum-tags-show", p.showTagH))
	// htm.GET("/articles", p.Layout.Application("forum-articles-index", p.indexArticlesH))
	// htm.GET("/articles/show/:id", p.Layout.Application("forum-articles-show", p.showArticleH))
	// htm.GET("/comments", p.Layout.Application("forum-comments-index", p.indexCommentsH))
	//
	// api := p.Router.Group("/api/forum", p.Layout.MustSignInMiddleware)
	// api.GET("/articles", p.Layout.JSON(p.indexArticles))
	// api.GET("/articles/:id", p.Layout.JSON(p.showArticle))
	// api.POST("/articles", p.Layout.JSON(p.createArticle))
	// api.POST("/articles/:id", p.canEditArticle, p.Layout.JSON(p.updateArticle))
	// api.DELETE("/articles/:id", p.canEditArticle, p.Layout.JSON(p.destroyArticle))
	// api.GET("/tags", p.Layout.JSON(p.indexTags))
	// api.GET("/tags/:id", p.Layout.JSON(p.showTag))
	// api.POST("/tags", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.createTag))
	// api.POST("/tags/:id", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.updateTag))
	// api.DELETE("/tags/:id", p.Layout.MustAdminMiddleware, p.Layout.JSON(p.destroyTag))
	// api.GET("/comments", p.Layout.JSON(p.indexComments))
	// api.GET("/comments/:id", p.Layout.JSON(p.showComment))
	// api.POST("/comments", p.Layout.JSON(p.createComment))
	// api.POST("/comments/:id", p.canEditComment, p.Layout.JSON(p.updateComment))
	// api.DELETE("/comments/:id", p.canEditComment, p.Layout.JSON(p.destroyComment))

	return nil
}

func init() {
	web.Register(&Plugin{})
}
