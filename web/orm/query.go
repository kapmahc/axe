package orm

import (
	log "github.com/sirupsen/logrus"
)

var queries = make(map[string]string, 0)

// Q get sql script by name. [db/queries.ini]
func Q(n string) string {
	q := queries[n]
	log.Debug("query ", n, " => ", q)
	return q
}
