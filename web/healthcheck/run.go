package healthcheck

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

// Run run
func Run() {
	for {
		do()
		time.Sleep(5 * time.Minute)
	}
}

func do() {
	for _, h := range handlers {
		log.Info("health check ", h.Name())
		val, err := h.Do()
		if err != nil {
			log.Error(err)
		}
		buf, err := json.Marshal(val)
		if err != nil {
			log.Error(err)
			return
		}
		log.Info(string(buf))
	}
}
