package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/rabbitmq"
	"github.com/dmitrykharchenko95/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/sirupsen/logrus"
)

type Scheduler struct {
	Producer *rabbitmq.Producer
	Storage  storage.Storage
	Logger   *logrus.Logger
	Period   time.Duration
}

func NewScheduler(prod *rabbitmq.Producer,
	store storage.Storage,
	logg *logrus.Logger,
	period time.Duration) *Scheduler {
	return &Scheduler{
		Producer: prod,
		Storage:  store,
		Logger:   logg,
		Period:   period,
	}
}

func (s *Scheduler) SearchEventForNotification(ctx context.Context) error {
	if err := s.Producer.Start(); err != nil {
		return fmt.Errorf("scheduler start error: %w", err)
	}

	defer ctx.Done()

	s.Logger.Info("Scheduler started...")

	ticker := time.NewTicker(s.Period)
	defer ticker.Stop()

	if err := s.Storage.Connect(ctx); err != nil {
		s.Logger.Errorf("storage connect fail:%v", err)
		return err
	}

	for {
		select {
		case <-ticker.C:
			err := s.Storage.DeleteOldEvents(ctx)
			if err != nil {
				s.Logger.Errorf("can not delete old events: %v", err)
				continue
			}

			events, err := s.Storage.GetEventForNotification(ctx, s.Period)
			if err != nil {
				s.Logger.Errorf("can not get events for notification: %v", err)
				continue
			}

			for _, e := range events {
				mes, err := parseEventForMessage(e)
				if err != nil {
					s.Logger.Errorf("can not parse event: %v", err)
					continue
				}

				if err = s.Producer.Publish(ctx, mes); err != nil {
					s.Logger.Errorf("can not published message: %v", err)
				}
			}
		case <-ctx.Done():
			s.Producer.Stop()
			return nil
		}
	}
}

func parseEventForMessage(e storage.Event) (rabbitmq.Message, error) {
	n := rabbitmq.Notification{
		ID:        e.ID,
		UserID:    e.UserID,
		Title:     e.Title,
		StartDate: e.StartDate,
	}

	data, err := json.Marshal(n)
	if err != nil {
		return rabbitmq.Message{}, err
	}
	return rabbitmq.Message{
		ContentType: "application/json",
		Body:        data,
	}, nil
}
