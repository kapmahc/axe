package nut

import (
	"net/http"

	"github.com/kapmahc/axe/web"
)

// MustSignInMiddleware currend user middleware
type MustSignInMiddleware struct {
	Wrapper *web.Wrapper `inject:""`
	I18n    *web.I18n    `inject:""`
}

func (p *MustSignInMiddleware) ServeHTTP(wrt http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	ctx := p.Wrapper.Context(wrt, req)
	if it, ok := ctx.Get(CurrentUser).(*User); ok && it.IsConfirm() && !it.IsLock() {
		next(wrt, req)
		return
	}
	ctx.Abort(http.StatusForbidden, p.I18n.E(ctx.Get(web.LOCALE).(string), "errors.forbidden"))
}
