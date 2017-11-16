package reading

import (
	"path"

	"github.com/facebookgo/inject"
	"github.com/go-pg/pg"
	"github.com/kapmahc/axe/plugins/nut"
	"github.com/kapmahc/axe/web"
	"github.com/kapmahc/stardict"
	log "github.com/sirupsen/logrus"
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

	dictionaries []*stardict.Dictionary
}

// Init init beans
func (p *Plugin) Init(*inject.Graph) error {
	dic := path.Join("tmp", "dic")
	log.Info("open stardict directories ", dic)
	if items, err := stardict.Open(dic); err == nil {
		p.dictionaries = items
	} else {
		return err
	}
	return nil
}

// Mount register
func (p *Plugin) Mount() error {
	return nil
}

func init() {
	web.Register(&Plugin{})

}
