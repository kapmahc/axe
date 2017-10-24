package nut

// Mount register
func (p *HomePlugin) Mount() error {
	p.Router.GET("/site/info", p.Layout.JSON(p.getAPISiteInfo))
	p.Router.POST("/install", p.Layout.JSON(p.postInstall))

	return nil
}
