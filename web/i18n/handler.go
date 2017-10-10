package i18n

import (
	"context"
	"net/http"

	"github.com/kapmahc/axe/web"
	"golang.org/x/text/language"
)

// LOCALE locale key
const LOCALE = web.K("locale")

// Middleware parse locales
func Middleware(wrt http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	name := string(LOCALE)
	matcher := language.NewMatcher(_languages)
	lang, written := detect(req, name)
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
	ctx := context.WithValue(req.Context(), LOCALE, lang)
	ctx = context.WithValue(ctx, web.PAYLOAD, web.H{
		"locale":    lang,
		"languages": _languages,
	})
	next(wrt, req.WithContext(ctx))
}

func detect(r *http.Request, k string) (string, bool) {
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

// Languages languages
// func Languages() ([]string, error) {
// 	rows, err := orm.DB().Query(orm.Q("i18n.languages"))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var items []string
// 	for rows.Next() {
// 		var name string
// 		if err = rows.Scan(&name); err != nil {
// 			return nil, err
// 		}
// 		items = append(items, name)
// 	}
// 	return items, nil
// }
