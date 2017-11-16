package web

// http://www.gnu.org/software/gettext/manual/gettext.html#Language-Codes
import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/go-ini/ini"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

// NewI18n create i18n
func NewI18n(path string, db *pg.DB) (*I18n, error) {
	it := I18n{
		db:    db,
		items: make(map[string]string),
	}
	if err := it.loadFromFileSystem(path); err != nil {
		return nil, err
	}
	return &it, nil
}

// Locale locale
type Locale struct {
	tableName struct{}  `sql:"locales"`
	ID        uint      `json:"id"`
	Lang      string    `json:"lang"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// I18n i18n
type I18n struct {
	db    *pg.DB
	items map[string]string
}

// Languages language tags
func (p *I18n) Languages() ([]string, error) {
	var langs []string
	if err := p.db.Model(&Locale{}).ColumnExpr("DISTINCT lang").Select(&langs); err != nil {
		return nil, err
	}
	if len(langs) == 0 {
		langs = append(langs, language.AmericanEnglish.String())
	}
	return langs, nil
}

func (p *I18n) loadFromFileSystem(dir string) error {
	const ext = ".ini"
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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
		log.Info("find locale ", tag)
		lang := tag.String()

		cfg, err := ini.Load(path)
		if err != nil {
			return err
		}

		for _, sec := range cfg.Sections() {
			z := sec.Name()
			for k, v := range sec.KeysHash() {
				p.items[lang+"."+z+"."+k] = v
			}
		}

		return nil
	})
}

// Set set
func (p *I18n) Set(tx *pg.Tx, lang, code, message string) error {
	var it Locale
	now := time.Now()
	err := tx.Model(&it).
		Column("id").
		Where("lang = ? AND code = ?", lang, code).
		Limit(1).Select()
	if err == nil {
		it.UpdatedAt = now
		it.Message = message
		_, err = tx.Model(&it).Column("message", "updated_at").Update()
	} else if err == pg.ErrNoRows {
		err = tx.Insert(&Locale{
			Lang:      lang,
			Code:      code,
			Message:   message,
			UpdatedAt: now,
		})
	}

	if err == nil {
		p.items[lang+"."+code] = message
	}
	return err
}

// H html
func (p *I18n) H(lang, code string, obj interface{}) (string, error) {
	msg, err := p.get(lang, code)
	if err != nil {
		return "", err
	}
	tpl, err := template.New("").Parse(msg)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, obj)
	return buf.String(), err
}

//E error
func (p *I18n) E(lang, code string, args ...interface{}) error {
	msg, err := p.get(lang, code)
	if err != nil {
		return err
	}
	return fmt.Errorf(msg, args...)
}

//T text
func (p *I18n) T(lang, code string, args ...interface{}) string {
	msg, err := p.get(lang, code)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf(msg, args...)
}

func (p *I18n) get(lang, code string) (string, error) {
	var it Locale
	if err := p.db.Model(&it).
		Column("message").
		Where("lang = ? AND code = ?", lang, code).
		Limit(1).Select(); err == nil {
		return it.Message, nil
	}
	key := lang + "." + code
	if msg, ok := p.items[key]; ok {
		return msg, nil
	}
	return "", errors.New(key)
}
