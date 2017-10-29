package nut

import "github.com/gin-gonic/gin"

func (p *AdminPlugin) indexLeaveWords(l string, c *gin.Context) (interface{}, error) {
	var items []LeaveWord
	err := p.DB.Model(&items).
		Order("created_at DESC").Select()
	return items, err
}
func (p *AdminPlugin) destroyLeaveWord(l string, c *gin.Context) (interface{}, error) {
	_, err := p.DB.Model(&LeaveWord{}).Where("id = ?", c.Param("id")).Delete()
	return gin.H{}, err
}
