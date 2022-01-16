package rabbitmq

import (
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	tag          string
	done         chan error
	uri          string
	exchangeName string
	exchangeType string
	queue        string
	bindingKey   string
}

type Handle func(delCh <-chan amqp.Delivery, done chan error)

func NewConsumer(port int,
	host,
	user,
	password,
	exchangeName,
	exchangeType,
	queueName,
	bindingKey,
	ctag string) *Consumer {
	uri := url.URL{
		Scheme: "amqp",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, strconv.Itoa(port)),
	}

	return &Consumer{
		tag:          ctag,
		done:         make(chan error),
		uri:          uri.String(),
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		queue:        queueName,
		bindingKey:   bindingKey,
	}
}

func (c *Consumer) Consume(handle Handle) error {
	var err error

	c.conn, err = amqp.Dial(c.uri)
	if err != nil {
		return fmt.Errorf("can not dial: %w", err)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("can not make channel: %w", err)
	}

	if err = c.channel.ExchangeDeclare(
		c.exchangeName,
		c.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("can not declare exchange : %w", err)
	}

	queue, err := c.channel.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("can not declare queue : %w", err)
	}

	if err = c.channel.QueueBind(
		queue.Name,
		c.bindingKey,
		c.exchangeName,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("can not bind queue: %w", err)
	}

	deliveries, err := c.channel.Consume(
		queue.Name,
		c.tag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("can not consume queue: %w", err)
	}

	handle(deliveries, c.done)

	return nil
}

func (c *Consumer) Shutdown() error {
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %w", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %w", err)
	}

	return <-c.done
}
