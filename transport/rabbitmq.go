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
	consumeC          <-chan amqp.Delivery
}

func (mq *RabbitMQ) Request(ctx context.Context, req interface{},
	resp interface{}) error {
	err := mq.init()
	if err != nil {
		return err
	}

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

func (mq *RabbitMQ) init() error {
	var err error
	log := func(format string, v ...interface{}) {
		mq.l.Printf("init", format, v...)
	}

	mq.initOnce.Do(func() {
		mq.setupLogger()

		mq.responseCh = make(chan []byte)
		consumeQueue, err1 := mq.Channel.QueueDeclare("", false, true, false, false,
			nil)
		if err1 != nil {
			err = fmt.Errorf("Channel.QueueDeclare(...): %w", err)
			return
		}

		log("queue: %s, exchange: %s, key: %s", consumeQueue.Name,
			mq.ConsumeFrom, mq.ConsumeKey)
		err1 = mq.Channel.QueueBind(consumeQueue.Name, mq.ConsumeKey, mq.ConsumeFrom, false,
			nil)
		if err1 != nil {
			err = fmt.Errorf("Channel.QueueBind(%s, %s, %s, ...): %w", consumeQueue.Name,
				mq.ConsumeKey, mq.ConsumeFrom, err1)
			return
		}
		ch, err1 := mq.Channel.Consume(consumeQueue.Name, "", true, false, false,
			false, nil)
		if err1 != nil {
			err = fmt.Errorf("Channel.Consume(%s, ...): %w", mq.ConsumeFrom, err)
			return
		}
		mq.consumeC = ch
	})

	return err
}

// ProcessEvents reads and handles events from RabbitMQ.Channel.
// This method is blocking and returns only if transport initialization failed,
// or after RabbitMQ.Channel closed. Returned value is always non nil error.
func (mq *RabbitMQ) ProcessEvents() error {
	err := mq.init()
	if err != nil {
		return err
	}

	log := func(format string, v ...interface{}) {
		mq.l.Printf("ProcessEvents", format, v...)
	}

	for d := range mq.consumeC {
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
