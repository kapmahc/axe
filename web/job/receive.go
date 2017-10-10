package job

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"

	log "github.com/sirupsen/logrus"
)

func receive(consumer string) error {
	log.Info("waiting for messages, to exit press CTRL+C")
	return open(func(ch *amqp.Channel) error {
		if err := ch.Qos(1, 0, false); err != nil {
			return err
		}
		qu, err := ch.QueueDeclare(_queue, true, false, false, false, nil)
		if err != nil {
			return err
		}
		msgs, err := ch.Consume(qu.Name, consumer, false, false, false, false, nil)
		if err != nil {
			return err
		}
		for d := range msgs {
			d.Ack(false)
			log.Info("receive message ", d.MessageId, " @ ", d.Type)
			now := time.Now()
			hnd, ok := handlers[d.Type]
			if !ok {
				return fmt.Errorf("unknown message type %s", d.Type)
			}
			if err := hnd(d.MessageId, d.Body); err != nil {
				return err
			}
			log.Info("done", d.MessageId, time.Now().Sub(now))
		}
		return nil
	})
}
