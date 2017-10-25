package nut

// Mount register
func (p *HomePlugin) Mount() error {
	htm := p.Router
	htm.GET("/", p.Layout.Application("nut/home", p.getHome))

	api := p.Router.Group("/api")
	api.GET("/site/info", p.Layout.JSON(p.getAPISiteInfo))
	api.POST("/install", p.Layout.JSON(p.postInstall))

	return nil
}
