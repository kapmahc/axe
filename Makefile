dist=dist
pkg=github.com/kapmahc/axe/web

VERSION=`git rev-parse --short HEAD`
BUILD_TIME=`date -R`
AUTHOR_NAME=`git config --get user.name`
AUTHOR_EMAIL=`git config --get user.email`
COPYRIGHT=`head -n 1 LICENSE`
USAGE=`sed -n '3p' README.md`

build: backend frontend
	cd $(dist) && tar cfJ ../$(dist).tar.xz *


backend:
	go build -ldflags "-s -w -X ${pkg}.Version=${VERSION} -X '${pkg}.BuildTime=${BUILD_TIME}' -X '${pkg}.AuthorName=${AUTHOR_NAME}' -X ${pkg}.AuthorEmail=${AUTHOR_EMAIL} -X '${pkg}.Copyright=${COPYRIGHT}' -X '${pkg}.Usage=${USAGE}'" -o ${dist}/axe main.go
	-cp -r locales templates $(dist)/


frontend:
	cd desktop && npm run build
	-cp -r desktop/package.json desktop/package-lock.json desktop/.next $(dist)/dashboard


clean:
	-rm -r $(dist) $(dist).tar.xz
	-rm -r desktop/.next
