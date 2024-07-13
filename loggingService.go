package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type loggingService struct {
	next PriceFetcher
}

func NewLogginService(next PriceFetcher) PriceFetcher {
	return &loggingService{
		next: next,
	}
}

func (s *loggingService) FetchPrice(ctx context.Context, ticker string) (price float64, err error) {
	defer func(begin time.Time) {
		logrus.WithFields(logrus.Fields{
			"requestId": ctx.Value(Key("requestId")),
			"took":      time.Since(begin),
			"err":       err,
			"price":     price,
		}).Info("fetch price")
	}(time.Now())
	return s.next.FetchPrice(ctx, ticker)
}
