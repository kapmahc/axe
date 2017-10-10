package orm

import (
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
)

// Load read queries
func Load() error {
	cfg, err := ini.Load(filepath.Join("db", "mapper.ini"))
	if err != nil {
		return err
	}

	for _, sec := range cfg.Sections() {
		for k, v := range sec.KeysHash() {
			queries[sec.Name()+"."+k] = v
		}
	}
	return nil
}

func versions() ([]string, error) {
	var items []string
	if err := filepath.Walk(filepath.Join("db", "migrations"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := filepath.Ext(info.Name())
		if !info.IsDir() {
			log.Warn("ingnore file ", path)
			return nil
		}
		log.Info("find migration ", name)
		items = append(items, name)
		return nil
	}); err != nil {
		return nil, err
	}
	return items, nil
}
