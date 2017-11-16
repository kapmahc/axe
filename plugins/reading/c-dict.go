package reading

import (
	"github.com/gin-gonic/gin"
)

type fmDict struct {
	Keywords string `json:"keywords" binding:"required,max=255"`
}

func (p *Plugin) postDict(c *gin.Context) (interface{}, error) {
	var fm fmDict

	if err := c.BindJSON(&fm); err != nil {
		return nil, err
	}
	rst := gin.H{}
	for _, dic := range p.dictionaries {
		for _, sen := range dic.Translate(fm.Keywords) {
			var items []gin.H
			for _, pat := range sen.Parts {
				items = append(items, gin.H{"type": pat.Type, "body": string(pat.Data)})
			}
			rst[dic.GetBookName()] = items
		}
	}
	return rst, nil
}
