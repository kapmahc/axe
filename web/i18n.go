package web

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-ini/ini"
	"github.com/go-pg/pg"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"golang.org/x/text/language"
)

// NewI18n create i18n
func NewI18n(path string, db *pg.DB) (*I18n, error) {
	it := I18n{
		items: make(map[string]string),
		langs: make([]language.Tag, 0),
	}
	if err := it.loadFromFileSystem(path); err != nil {
		return nil, err
	}
	if err := it.loadFromDb(db); err != nil {
		return nil, err
	}
	return &it, nil
}

// LOCALE locale context key
const LOCALE = "locale"

// Locale locale
type Locale struct {
	tableName struct{} `sql:"locales"`
	ID        uint
	Lang      string
	Code      string
	Message   string
	Updated   time.Time
	Created   time.Time
}

// I18n i18n
type I18n struct {
	items map[string]string
	langs []language.Tag
}

// Languages language tags
func (p *I18n) Languages() []language.Tag {
	return p.langs[:]
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
		p.langs = append(p.langs, tag)
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

func (p *I18n) loadFromDb(db *pg.DB) error {
	var items []Locale
	if err := db.Model(&items).
		Column("lang", "code", "message").
		Select(); err != nil {
		return err
	}
	for _, it := range items {
		p.items[it.Lang+"."+it.Code] = it.Message
	}
	return nil
}

// Set set
func (p *I18n) Set(tx *pg.Tx, lang, code, message string) error {
	var it Locale
	now := time.Now()
	err := tx.Model(&it).Column("id").Where("lang = ? AND code = ?", lang, code).Select()
	if err == nil {
		it.Updated = now
		it.Message = message
		_, err = tx.Model(&it).Column("message").Update()
	} else if err == pg.ErrNoRows {
		err = tx.Insert(&Locale{
			Lang:    lang,
			Code:    code,
			Message: message,
			Updated: now,
			Created: now,
		})
	}

	if err == nil {
		p.items[lang+"."+code] = message
	}
	return err
}

// H html
func (p *I18n) H(lang, code string, obj interface{}) (string, error) {
	k := lang + "." + code
	msg, ok := p.items[k]
	if !ok {
		return k, nil
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
	k := lang + "." + code
	msg, ok := p.items[k]
	if !ok {
		return errors.New(k)
	}
	return fmt.Errorf(msg, args...)
}

//T text
func (p *I18n) T(lang, code string, args ...interface{}) string {
	k := lang + "." + code
	msg, ok := p.items[k]
	if !ok {
		return k
	}
	return fmt.Sprintf(msg, args...)
}

// Middleware parse locales
func (p *I18n) Middleware() negroni.HandlerFunc {
	name := string(LOCALE)
	matcher := language.NewMatcher(p.langs)
	return func(wrt http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

		lang, written := p.detect(req, name)
		tag, _, _ := matcher.Match(language.Make(lang))
		if lang != tag.String() {
			written = true
			lang = tag.String()
		}
		if written {
			http.SetCookie(wrt, &http.Cookie{
				Name:     name,
				Value:    lang,
				MaxAge:   1<<32 - 1,
				Secure:   false,
				HttpOnly: false,
			})
		}
		ctx := context.WithValue(req.Context(), K(LOCALE), lang)
		next(wrt, req.WithContext(ctx))
	}
}

func (p *I18n) detect(r *http.Request, k string) (string, bool) {
	// 1. Check URL arguments.
	if lang := r.URL.Query().Get(k); lang != "" {
		return lang, true
	}

	// 2. Get language information from cookies.
	if ck, er := r.Cookie(k); er == nil {
		return ck.Value, false
	}

	// 3. Get language information from 'Accept-Language'.
	return r.Header.Get("Accept-Language"), true
}
