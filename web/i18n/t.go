package i18n

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
)

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
