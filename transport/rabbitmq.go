package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/streadway/amqp"

	"github.com/dinalt/jango"
	"github.com/dinalt/jango/logging"
)

var _ jango.Transport = (*RabbitMQ)(nil)
var ErrChanClosed = fmt.Errorf("channel closed")

type RabbitMQ struct {
	Channel           *amqp.Channel
	PublishTo         string
	PublishKey        string
	Logger            logging.Logger
	l                 *logging.StructLogger
	ConsumeFrom       string
	ConsumeKey        string
	pendigTransaction string
	responseCh        chan []byte
	initOnce          sync.Once
}

func (mq *RabbitMQ) Request(ctx context.Context, req interface{},
	resp interface{}) error {
	mq.init()

	log := func(format string, v ...interface{}) {
		mq.l.Printf("Request", format, v...)
	}

	transactional, _ := req.(interface{ Transaction() string })
	if transactional == nil {
		panic("req does not have Transaction getter")
	}

	mq.pendigTransaction = transactional.Transaction()
	if mq.pendigTransaction == "" {
		panic("transaction is empty")
	}

	b, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("json.Marshal(...): %w", err)
	}

	p := amqp.Publishing{
		ContentType:   "application/json",
		Body:          b,
		CorrelationId: mq.pendigTransaction,
	}

	log("OUTCOMING: %s", b)
	err = mq.Channel.Publish(mq.PublishTo, mq.PublishKey, false, false, p)
	if err != nil {
		return fmt.Errorf("Channel.Publish(%s, %s, ...): %w", mq.PublishTo,
			mq.PublishKey, err)
	}

	select {
	case mqResp := <-mq.responseCh:
		mq.pendigTransaction = ""
		err := json.Unmarshal(mqResp, &resp)
		if err != nil {
			return fmt.Errorf("json.Unmarshal(...): %w", err)
		}
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (mq *RabbitMQ) init() {
	mq.initOnce.Do(mq.setupLogger)
}

func (mq *RabbitMQ) Serve() error {
	mq.init()

	log := func(format string, v ...interface{}) {
		mq.l.Printf("Serve", format, v...)
	}

	mq.responseCh = make(chan []byte)
	consumeQueue, err := mq.Channel.QueueDeclare("", false, false, false, false,
		nil)
	if err != nil {
		return fmt.Errorf("Channel.QueueDeclare(...): %w", err)
	}

	log("queue: %s, exchange: %s, key: %s", consumeQueue.Name,
		mq.ConsumeFrom, mq.ConsumeKey)
	err = mq.Channel.QueueBind(consumeQueue.Name, mq.ConsumeKey, mq.ConsumeFrom, false,
		nil)
	if err != nil {
		return fmt.Errorf("Channel.QueueBind(%s, %s, %s, ...): %w", consumeQueue.Name,
			mq.ConsumeKey, mq.ConsumeFrom, err)
	}
	ch, err := mq.Channel.Consume(consumeQueue.Name, "", true, false, false,
		false, nil)
	if err != nil {
		return fmt.Errorf("Channel.Consume(%s, ...): %w", mq.ConsumeFrom, err)
	}

	for d := range ch {
		log("INCOMING: %s", d.Body)
		if d.CorrelationId != mq.pendigTransaction {
			continue
		}
		mq.responseCh <- d.Body
	}

	return ErrChanClosed
}

func (mq *RabbitMQ) setupLogger() {
	if mq.Logger == nil {
		mq.Logger = logging.NoopLogger{}
	}
	mq.l = logging.NewStructLogger(RabbitMQ{}, mq.Logger, "github.com/dinalt/")
}
