package nut

// Mount register
func (p *HomePlugin) Mount() error {
	rt := p.Router.Group("/api")
	rt.GET("/site/info", p.Layout.JSON(p.getAPISiteInfo))
	rt.POST("/install", p.Layout.JSON(p.postInstall))

	return nil
}
