package main

import (
	"log"

	_ "github.com/kapmahc/axe/plugins/erp"
	_ "github.com/kapmahc/axe/plugins/forum"
	_ "github.com/kapmahc/axe/plugins/mall"
	"github.com/kapmahc/axe/plugins/nut"
	_ "github.com/kapmahc/axe/plugins/ops/mail"
	_ "github.com/kapmahc/axe/plugins/ops/vpn"
	_ "github.com/kapmahc/axe/plugins/pos"
	_ "github.com/kapmahc/axe/plugins/reading"
	_ "github.com/kapmahc/axe/plugins/survey"
	_ "github.com/lib/pq"
)

func main() {
	if err := nut.Main(); err != nil {
		log.Fatal(err)
	}
}
