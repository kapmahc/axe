package i18n

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"html/template"

	"github.com/kapmahc/axe/web/orm"
)

// Languages languages
func Languages(db *sql.DB) ([]string, error) {
	rows, err := db.Query(orm.Q("i18n.languages"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	return items, nil
}

// H html
func H(lang, code string, obj interface{}) (string, error) {
	k := key(lang, code)
	msg, ok := _locales[k]
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
func E(lang, code string, args ...interface{}) error {
	k := key(lang, code)
	msg, ok := _locales[k]
	if !ok {
		return errors.New(k)
	}
	return fmt.Errorf(msg, args...)
}

//T text
func T(lang, code string, args ...interface{}) string {
	k := key(lang, code)
	msg, ok := _locales[k]
	if !ok {
		return k
	}
	return fmt.Sprintf(msg, args...)
}
