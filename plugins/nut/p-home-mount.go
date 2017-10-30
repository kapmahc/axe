package nut

// Mount register
func (p *HomePlugin) Mount() error {
	htm := p.Router
	htm.GET("/", p.Layout.Application("nut-home.html", p.getHome))

	api := p.Router.Group("/api")
	api.POST("/token", p.Layout.MustSignInMiddleware, p.Layout.JSON(p.postToken))
	api.GET("/site/info", p.Layout.JSON(p.getSiteInfo))
	api.POST("/install", p.Layout.JSON(p.postInstall))
	api.POST("/leave-words", p.Layout.JSON(p.createLeaveWord))

	return nil
}
