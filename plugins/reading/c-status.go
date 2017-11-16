package reading

import (
	"github.com/gin-gonic/gin"
)

func (p *Plugin) getStatus(c *gin.Context) (interface{}, error) {
	data := gin.H{}
	bc, err := p.DB.Model(&Book{}).Count()
	if err != nil {
		return nil, err
	}
	data["book"] = gin.H{
		"count": bc,
	}

	dict := gin.H{}
	for _, dic := range p.dictionaries {
		dict[dic.GetBookName()] = dic.GetWordCount()
	}
	data["dict"] = dict

	return data, nil

}
