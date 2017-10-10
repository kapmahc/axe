package job

import (
	log "github.com/sirupsen/logrus"
)

// Handler handler
type Handler func(id string, body []byte) error

var handlers = make(map[string]Handler)

// Register register handler
func Register(queue string, hnd Handler) {
	if _, ok := handlers[queue]; ok {
		log.Warn("handler for queue ", queue, " already exists, will override it")
	}
	handlers[queue] = hnd
}
