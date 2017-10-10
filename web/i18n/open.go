package i18n

import (
	"os"
	"path/filepath"

	"github.com/go-ini/ini"
	"github.com/kapmahc/axe/web/orm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

var (
	_languages []language.Tag
)

// Open load locales from database, filesystem
func Open(dir string) error {
	if err := loadFromFileSystem(dir); err != nil {
		return err
	}
	if err := loadFromDb(); err != nil {
		return err
	}
	return nil
}

func loadFromFileSystem(dir string) error {
	return filepath.Walk("locales", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := info.Name()
		if info.IsDir() || filepath.Ext(name) != ext {
			return err
		}
		tag, err := language.Parse(name[:len(name)-len(ext)])
		if err != nil {
			return err
		}
		lang := tag.String()
		log.Info("find locale ", lang)
		_languages = append(_languages, tag)

		cfg, err := ini.Load(path)
		if err != nil {
			return err
		}

		for _, sec := range cfg.Sections() {
			z := sec.Name()
			for k, v := range sec.KeysHash() {
				_locales[lang+sep+z+sep+k] = v
			}
		}

		return nil
	})
}

func loadFromDb() error {
	rows, err := orm.DB().Query(orm.Q("i18n.locales"))
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var lang, code, message string
		if err := rows.Scan(&lang, &code, &message); err != nil {
			return err
		}
		_locales[lang+sep+code] = code
	}
	return nil
}
