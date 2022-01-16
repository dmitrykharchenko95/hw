package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/streadway/amqp"
)

type Producer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	reConn       chan *amqp.Error
	uri          string
	exchangeName string
	exchangeType string
	queue        string
	bindingKey   string
}

func NewProducer(port int, host, user, password, exchangeName, exchangeType, queueName, bindingKey string) *Producer {
	uri := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, strconv.Itoa(port)),
	}

	return &Producer{
		uri:          uri.String(),
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		queue:        queueName,
		bindingKey:   bindingKey,
		reConn:       make(chan *amqp.Error),
	}
}

func (p *Producer) Start() error {
	var err error
	if err = p.connect(); err != nil {
		return fmt.Errorf("producer connect: %w", err)
	}
	if err = p.channel.ExchangeDeclare(
		p.exchangeName,
		p.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("exchange declare: %w", err)
	}
	_, err = p.channel.QueueDeclare(
		p.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue declare: %w", err)
	}

	if err = p.channel.QueueBind(
		p.queue,
		p.bindingKey,
		p.exchangeName,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("bind queue: %w", err)
	}

	return nil
}

func (p *Producer) Publish(ctx context.Context, m Message) error {
	select {
	case <-p.reConn:
		err := p.reConnect(ctx)
		if err != nil {
			return fmt.Errorf("reconnecting error: %w", err)
		}
	default:
	}

	return p.channel.Publish(
		p.exchangeName,
		p.bindingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     m.ContentType,
			ContentEncoding: "",
			Body:            m.Body,
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		},
	)
}

func (p *Producer) connect() error {
	var err error

	p.conn, err = amqp.Dial(p.uri)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}

	p.channel, err = p.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %w", err)
	}
	p.conn.NotifyClose(p.reConn)

	return nil
}

func (p *Producer) reConnect(ctx context.Context) error {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, ctx)
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return fmt.Errorf("stop reconnecting")
		}
		time.Sleep(d)
		if err := p.connect(); err != nil {
			log.Printf("could not connect in reconnect call: %+v", err)
			continue
		}

		return nil
	}
}

func (p *Producer) Stop() {
	err := p.channel.Close()
	if err != nil {
		log.Printf("can not close producer chanal:%v", err)
	}
	err = p.conn.Close()
	if err != nil {
		log.Printf("can not close producer connect:%v", err)
	}
}
