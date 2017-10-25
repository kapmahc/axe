package nut

// Mount register
func (p *HomePlugin) Mount() error {
	api := p.Router.Group("/api")
	api.GET("/site/info", p.Layout.JSON(p.getAPISiteInfo))
	api.POST("/install", p.Layout.JSON(p.postInstall))

	return nil
}
