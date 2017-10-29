package nut

import "github.com/gin-gonic/gin"

func (p *AdminPlugin) indexUsers(l string, c *gin.Context) (interface{}, error) {
	var items []User
	err := p.DB.Model(&items).
		Column("id", "name", "email",
			"sign_in_count",
			"last_sign_in_at", "last_sign_in_ip",
			"current_sign_in_at", "current_sign_in_ip",
		).
		Order("updated_at DESC").Select()
	return items, err
}
