package job

import "github.com/streadway/amqp"

var (
	_url   string
	_queue string
)

// Open open
func Open(url string, queue string) {
	_url = url
	_queue = queue
}

func open(f func(*amqp.Channel) error) error {
	conn, err := amqp.Dial(_url)
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return f(ch)
}
