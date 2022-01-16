package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Sender struct {
	Logger   *logrus.Logger
	Consumer *rabbitmq.Consumer
}

func NewSender(logger *logrus.Logger, consumer *rabbitmq.Consumer) *Sender {
	return &Sender{
		Logger:   logger,
		Consumer: consumer,
	}
}

func (s *Sender) ListenQueue(ctx context.Context) error {
	s.Logger.Info("Sender started...")

	defer func(Consumer *rabbitmq.Consumer) {
		err := Consumer.Shutdown()
		if err != nil {
			s.Logger.Warnf("consumer shutdown fail:%v", err)
		}
		ctx.Done()
	}(s.Consumer)

	if err := s.Consumer.Consume(handle); err != nil {
		s.Logger.Warnf("Listen queue fail:%v", err)
		return err
	}
	return nil
}

var handle rabbitmq.Handle = func(delCh <-chan amqp.Delivery, done chan error) {
	notif := &rabbitmq.Notification{}

	for d := range delCh {
		select {
		case err := <-done:
			log.Printf("done case: %v", err)
			return
		default:
			err := json.Unmarshal(d.Body, notif)
			if err != nil {
				log.Printf("notification unmarshal fail:%v", err)
				continue
			}
			fmt.Printf("New notification for user #%d\nEvent #%d: %v\nStart - %v",
				notif.UserID, notif.ID, notif.Title, notif.StartDate)
		}
	}
}
