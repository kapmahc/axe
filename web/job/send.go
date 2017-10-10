package job

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// Send send job
func Send(_type string, priority uint8, body interface{}) error {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(body)
	if err != nil {
		return err
	}
	return open(func(ch *amqp.Channel) error {
		qu, err := ch.QueueDeclare(_queue, true, false, false, false, nil)
		if err != nil {
			return err
		}

		return ch.Publish("", qu.Name, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.New().String(),
			Priority:     priority,
			Body:         buf.Bytes(),
			Timestamp:    time.Now(),
			Type:         _type,
		})
	})
}
